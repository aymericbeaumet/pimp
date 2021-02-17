package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"reflect"

	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v3"
)

// TODO: load from config file
var configStr = `
mappings:

  git :
    hub status -sb

  git st :
    hub status

  git co ... :
    hub checkout ...

  git ... :
    hub ...
`

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: prox cmd [argument]...")
		return
	}

	engine, err := NewEngineFromString(configStr)
	if err != nil {
		panic(err)
	}

	env, args := engine.Map(os.Environ(), os.Args[1:])
	fmt.Println(env, args)
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
	Mappings map[string]string
}

type Engine struct {
	mappings []*Mapping
}

type Mapping struct {
	from []string
	to   []string
	env  []string
}

func NewEngineFromString(configString string) (*Engine, error) {
	engine := &Engine{}

	var config Config
	if err := yaml.Unmarshal([]byte(configString), &config); err != nil {
		return nil, err
	}

	for fromString, toString := range config.Mappings {
		from, err := shellwords.Parse(fromString)
		if err != nil {
			panic(err)
		}
		env, to, err := shellwords.ParseWithEnvs(toString)
		if err != nil {
			panic(err)
		}
		engine.mappings = append(engine.mappings, &Mapping{
			from: from,
			to:   to,
			env:  env,
		})
	}

	return engine, nil
}

func (e *Engine) Map(env []string, args []string) ([]string, []string) {
	for _, mapping := range e.mappings {
		if mapping.from[len(mapping.from)-1] == "..." {
			from := mapping.from[:len(mapping.from)-1]
			if reflect.DeepEqual(args[:len(from)], from) {
				to := append(mapping.to[:len(mapping.to)-1], args[len(from):]...)
				return append(env, mapping.env...), to
			}
		}

		if reflect.DeepEqual(args, mapping.from) {
			return append(env, mapping.env...), mapping.to
		}
	}

	return env, args
}
