package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: pimp [option]... cmd [arg]...")
		return
	}

	engine, err := NewEngineFromHomeConfig()
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

	signalC := make(chan os.Signal, 32)
	signal.Notify(signalC)
	go func() {
		for signal := range signalC {
			_ = cmd.Process.Signal(signal)
		}
	}()

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

func NewEngineFromHomeConfig() (*Engine, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ".pimprc")
	return NewEngineFromConfig(configPath)
}

func NewEngineFromConfig(configPath string) (*Engine, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	return NewEngineFromReader(file)
}

func NewEngineFromReader(r io.Reader) (*Engine, error) {
	engine := &Engine{}

	var config Config
	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
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

				file, err := ioutil.TempFile("", "pimp")
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
