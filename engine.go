package main

import (
	"encoding/json"
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
	Sources  map[string]*Source `json:"sources"`
	Mappings []*Mapping         `json:"mappings"`

	// used to cache calls to the Executables() method
	executables []string
}

type Source struct {
	ID   string   `json:"id"`
	Env  []string `json:"env"`
	Args []string `json:"args"`
}

type Mapping struct {
	Pattern []string `json:"pattern"`
	Env     []string `json:"env"`
	Args    []string `json:"args"`
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

	engine.Sources = make(map[string]*Source, len(config.Sources))
	for id, raw := range config.Sources {
		env, args, err := parseEnvArgs(raw)
		if err != nil {
			return nil, err
		}
		engine.Sources[id] = &Source{
			ID:   id,
			Env:  env,
			Args: args,
		}
	}

	for _, mapping := range config.Mappings {
		for rawPattern, raw := range mapping {
			pattern, err := shellwords.Parse(rawPattern)
			if err != nil {
				return nil, err
			}

			env, args, err := parseEnvArgs(raw)
			if err != nil {
				return nil, err
			}

			engine.Mappings = append(engine.Mappings, &Mapping{
				Pattern: pattern,
				Env:     env,
				Args:    args,
			})
		}
	}

	return engine, nil
}

func (e *Engine) Map(env []string, args []string) ([]string, []string) {
	for _, mapping := range e.Mappings {
		if mapping.Pattern[len(mapping.Pattern)-1] == "..." {
			pattern := mapping.Pattern[:len(mapping.Pattern)-1]
			lim := len(pattern)
			if lim > len(args) {
				lim = len(args)
			}
			if reflect.DeepEqual(pattern, args[:lim]) {
				return append(env[:], mapping.Env...), append(mapping.Args[:], args[lim:]...)
			}
		}

		if reflect.DeepEqual(mapping.Pattern, args) {
			return append(env[:], mapping.Env...), mapping.Args
		}
	}

	return env, args
}

func (e *Engine) Dump(w io.Writer) error {
	return json.NewEncoder(w).Encode(e)
}

func (e *Engine) Executables() []string {
	if e.executables != nil {
		return e.executables
	}

	set := map[string]struct{}{}
	for _, m := range e.Mappings {
		set[m.Pattern[0]] = struct{}{}
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

func parseEnvArgs(input string) ([]string, []string, error) {
	const SHEBANG = "#!"

	// replace templates by placeholders
	templatesByPlaceholder := map[string]string{}
	input = templateRegexp.ReplaceAllStringFunc(input, func(template string) string {
		placeholder := fmt.Sprintf("___pimp_%d___", len(templatesByPlaceholder))
		templatesByPlaceholder[placeholder] = template
		return placeholder
	})

	var env, args []string
	var err error

	// multiline script (with shebang)
	if newLineIndex := strings.IndexRune(input, '\n'); newLineIndex > -1 {
		if !strings.HasPrefix(input, SHEBANG) {
			return nil, nil, errors.New("invalid shebang")
		}

		args, err = shellwords.Parse(input[len(SHEBANG):newLineIndex])
		if err != nil {
			return nil, nil, err
		}
		if len(args) == 0 || !strings.HasPrefix(args[0], "/") {
			return nil, nil, errors.New("shebang must be an absolute path")
		}

		file, err := ioutil.TempFile("", "pimp")
		if err != nil {
			return nil, nil, err
		}
		args = append(args, file.Name())

		if _, err := file.WriteString(input); err != nil {
			return nil, nil, err
		}
	} else {
		env, args, err = shellwords.ParseWithEnvs(input)
		if err != nil {
			return nil, nil, err
		}
	}

	// replace placeholders by templates
	for i, c := range args {
		if template, ok := templatesByPlaceholder[c]; ok {
			args[i] = template
		}
	}

	return env, args, nil
}
