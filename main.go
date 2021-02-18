package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v3"
)

// TODO: load from config file
var configStr = `
mappings:
  - git: hub status -sb

  - git a ...: hub add

  - git ci ...: hub commit

  - git co ...: hub checkout

  - git df ...: hub diff

  - git dfc ...: hub diff --cached

  - git l ...: hub log

  - git plps: |
      #!/bin/sh
      git pull && git push

  - git ps ...: hub push

  - git st ...: hub status

  - git _bash ...: |
      #!/bin/bash
      echo "Hello Bash! $@"

  - git _cat: |
      #!/bin/cat
      Hello Cat!

  - git _python ...: |
      #!/usr/bin/python
      import sys
      print("Hello Python!", sys.argv)

  - git _ruby ...: |
      #!/usr/bin/ruby
      puts "Hello Ruby!", ARGV

  - git ...: hub
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: prox [option]... cmd [arg]...")
		return
	}

	engine, err := NewEngineFromString(configStr)
	if err != nil {
		panic(err)
	}

	env, args := engine.Map(os.Environ(), os.Args[1:])
	cmd := exec.CommandContext(context.Background(), args[0], args[1:]...)
	cmd.Env = env
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	state, err := cmd.Process.Wait()
	if err != nil {
		panic(err)
	}
	os.Exit(state.ExitCode())
}

type Config struct {
	Mappings []map[string]string
}

type Engine struct {
	mappings []*Mapping
}

type Mapping struct {
	from []string
	to   []string
	env  []string
}

const SHEBANG = "#!"

func NewEngineFromString(configString string) (*Engine, error) {
	engine := &Engine{}

	var config Config
	if err := yaml.Unmarshal([]byte(configString), &config); err != nil {
		return nil, err
	}

	for _, mapping := range config.Mappings {
		for fromString, toString := range mapping {
			from, err := shellwords.Parse(fromString)
			if err != nil {
				panic(err)
			}

			var env []string
			var to []string
			// multiline script (with shebang)
			if newLineIndex := strings.IndexRune(toString, '\n'); newLineIndex > -1 {
				if !strings.HasPrefix(toString, SHEBANG) {
					panic("invalid shebang")
				}

				to, err = shellwords.Parse(toString[len(SHEBANG):newLineIndex])
				if err != nil {
					panic(err)
				}
				if len(to) == 0 || !strings.HasPrefix(to[0], "/") {
					panic("shebang must be an absolute path")
				}

				file, err := ioutil.TempFile("", "prox")
				if err != nil {
					panic(err)
				}
				to = append(to, file.Name())

				if _, err := file.WriteString(toString); err != nil {
					panic(err)
				}
			} else { // single line command
				env, to, err = shellwords.ParseWithEnvs(toString)
				if err != nil {
					panic(err)
				}
			}

			engine.mappings = append(engine.mappings, &Mapping{
				from: from,
				to:   to,
				env:  env,
			})
		}
	}

	return engine, nil
}

func (e *Engine) Map(env []string, args []string) ([]string, []string) {
	for _, mapping := range e.mappings {
		if mapping.from[len(mapping.from)-1] == "..." {
			from := mapping.from[:len(mapping.from)-1]
			if reflect.DeepEqual(from, args[:len(from)]) {
				return append(env[:], mapping.env...), append(mapping.to[:], args[len(from):]...)
			}
		}

		if reflect.DeepEqual(mapping.from, args) {
			return append(env[:], mapping.env...), mapping.to
		}
	}

	return env, args
}
