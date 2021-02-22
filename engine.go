package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Sources  map[string]string   `yaml:"sources"`
	Mappings []map[string]string `yaml:"mappings"`
}

type Engine struct {
	sources  map[string]*Source
	mappings []*Mapping

	// used to cache calls to the Executables() method
	executables []string
}

type Source struct {
	id string
	//
	env     []string
	command []string
}

type Mapping struct {
	pattern []string
	//
	env     []string
	command []string
}

func NewEngineFromFile(name string) (*Engine, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return NewEngineFromReader(file)
}

func NewEngineFromReader(r io.Reader) (*Engine, error) {
	engine := &Engine{}

	var config Config
	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
		return nil, err
	}

	engine.sources = make(map[string]*Source, len(config.Sources))
	for id, rawCommand := range config.Sources {
		env, command, err := parseToEnvCommand(rawCommand)
		if err != nil {
			return nil, err
		}
		engine.sources[id] = &Source{
			id:      id,
			env:     env,
			command: command,
		}
	}

	for _, mapping := range config.Mappings {
		for rawPattern, rawCommand := range mapping {
			pattern, err := shellwords.Parse(rawPattern)
			if err != nil {
				return nil, err
			}

			env, command, err := parseToEnvCommand(rawCommand)
			if err != nil {
				return nil, err
			}

			engine.mappings = append(engine.mappings, &Mapping{
				pattern: pattern,
				env:     env,
				command: command,
			})
		}
	}

	return engine, nil
}

func (e *Engine) Map(env []string, args []string) ([]string, []string) {
	for _, mapping := range e.mappings {
		if mapping.pattern[len(mapping.pattern)-1] == "..." {
			pattern := mapping.pattern[:len(mapping.pattern)-1]
			lim := len(pattern)
			if lim > len(args) {
				lim = len(args)
			}
			if reflect.DeepEqual(pattern, args[:lim]) {
				return append(env[:], mapping.env...), append(mapping.command[:], args[lim:]...)
			}
		}

		if reflect.DeepEqual(mapping.pattern, args) {
			return append(env[:], mapping.env...), mapping.command
		}
	}

	return env, args
}

func (e *Engine) Dump(w io.Writer) error {
	if _, err := fmt.Fprintf(w, "Sources:\n"); err != nil {
		return err
	}
	for _, source := range e.sources {
		if _, err := fmt.Fprintf(w, "  - %s: %#v\n", source.id, source.command); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintf(w, "\nMappings:\n"); err != nil {
		return err
	}
	for _, mapping := range e.mappings {
		if _, err := fmt.Fprintf(w, "  - %#v => %#v\n", mapping.pattern, mapping.command); err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) Executables() []string {
	if e.executables != nil {
		return e.executables
	}

	set := map[string]struct{}{}
	for _, m := range e.mappings {
		set[m.pattern[0]] = struct{}{}
	}

	out := make([]string, 0, len(set))
	for entry := range set {
		out = append(out, entry)
	}
	sort.Strings(out)

	e.executables = out
	return out
}

var templateRegexp = regexp.MustCompile(`{{[^}]+}}`)

func parseToEnvCommand(input string) ([]string, []string, error) {
	const SHEBANG = "#!"

	// replace templates by placeholders
	templatesByPlaceholder := map[string]string{}
	input = templateRegexp.ReplaceAllStringFunc(input, func(template string) string {
		placeholder := fmt.Sprintf("___pimp_%d___", len(templatesByPlaceholder))
		templatesByPlaceholder[placeholder] = template
		return placeholder
	})

	var env, command []string
	var err error

	// multiline script (with shebang)
	if newLineIndex := strings.IndexRune(input, '\n'); newLineIndex > -1 {
		if !strings.HasPrefix(input, SHEBANG) {
			return nil, nil, errors.New("invalid shebang")
		}

		command, err = shellwords.Parse(input[len(SHEBANG):newLineIndex])
		if err != nil {
			return nil, nil, err
		}
		if len(command) == 0 || !strings.HasPrefix(command[0], "/") {
			return nil, nil, errors.New("shebang must be an absolute path")
		}

		file, err := ioutil.TempFile("", "pimp")
		if err != nil {
			return nil, nil, err
		}
		command = append(command, file.Name())

		if _, err := file.WriteString(input); err != nil {
			return nil, nil, err
		}
	} else {
		env, command, err = shellwords.ParseWithEnvs(input)
		if err != nil {
			return nil, nil, err
		}
	}

	// replace placeholders by templates
	for i, c := range command {
		if template, ok := templatesByPlaceholder[c]; ok {
			command[i] = template
		}
	}

	return env, command, nil
}
