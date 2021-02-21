package main

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v2"
)

const SHEBANG = "#!"

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
				return nil, err
			}

			var env []string
			var to []string
			// multiline script (with shebang)
			if newLineIndex := strings.IndexRune(toString, '\n'); newLineIndex > -1 {
				if !strings.HasPrefix(toString, SHEBANG) {
					return nil, errors.New("invalid shebang")
				}

				to, err = shellwords.Parse(toString[len(SHEBANG):newLineIndex])
				if err != nil {
					return nil, err
				}
				if len(to) == 0 || !strings.HasPrefix(to[0], "/") {
					return nil, errors.New("shebang must be an absolute path")
				}

				file, err := ioutil.TempFile("", "pimp")
				if err != nil {
					return nil, err
				}
				to = append(to, file.Name())

				if _, err := file.WriteString(toString); err != nil {
					return nil, err
				}
			} else { // single line command
				env, to, err = shellwords.ParseWithEnvs(toString)
				if err != nil {
					return nil, err
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
			lim := len(from)
			if lim > len(args) {
				lim = len(args)
			}
			if reflect.DeepEqual(from, args[:lim]) {
				return append(env[:], mapping.env...), append(mapping.to[:], args[lim:]...)
			}
		}

		if reflect.DeepEqual(mapping.from, args) {
			return append(env[:], mapping.env...), mapping.to
		}
	}

	return env, args
}

func (e *Engine) Executables() []string {
	set := map[string]struct{}{}
	for _, m := range e.mappings {
		set[m.from[0]] = struct{}{}
	}

	out := make([]string, 0, len(set))
	for entry := range set {
		out = append(out, entry)
	}
	sort.Strings(out)

	return out
}
