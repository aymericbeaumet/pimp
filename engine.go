package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v2"
)

type Config yaml.MapSlice

type Engine struct {
	MappingsByArg0 map[string][]*Mapping `json:"mappings"`

	// used to cache calls to the Executables() method
	executables []string
}

type Mapping struct {
	Pattern []string          `json:"pattern"`
	Env     []string          `json:"env,omitempty"`
	Args    []string          `json:"args"`
	Files   map[string]string `json:"files,omitempty"`
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
	engine := &Engine{
		MappingsByArg0: map[string][]*Mapping{},
	}

	var config Config
	if err := yaml.NewDecoder(r).Decode(&config); err != nil {
		return nil, err
	}

	for _, item := range config {
		pattern, err := shellwords.Parse(item.Key.(string))
		if err != nil {
			return nil, err
		}

		env, args, files, err := parseEnvArgs(item.Value.(string))
		if err != nil {
			return nil, err
		}

		engine.MappingsByArg0[args[0]] = append(engine.MappingsByArg0[args[0]], &Mapping{
			Pattern: pattern,
			Env:     env,
			Args:    args,
			Files:   files,
		})
	}

	return engine, nil
}

func (e *Engine) Map(env []string, args []string) ([]string, []string, map[string]string) {
	mappings, ok := e.MappingsByArg0[args[0]]
	if !ok {
		return env, args, nil
	}

	for _, mapping := range mappings {
		if mapping.Pattern[len(mapping.Pattern)-1] == "..." {
			pattern := mapping.Pattern[:len(mapping.Pattern)-1]
			lim := len(pattern)
			if lim > len(args) {
				lim = len(args)
			}
			if reflect.DeepEqual(pattern, args[:lim]) {
				return append(env[:], mapping.Env...), append(mapping.Args[:], args[lim:]...), mapping.Files
			}
		}

		if reflect.DeepEqual(mapping.Pattern, args) {
			return append(env[:], mapping.Env...), mapping.Args, mapping.Files
		}
	}

	return env, args, nil
}

func (e *Engine) Dump(w io.Writer) error {
	return json.NewEncoder(w).Encode(e)
}

func (e *Engine) Executables() []string {
	if e.executables != nil {
		return e.executables
	}

	out := make([]string, 0, len(e.MappingsByArg0))
	for arg0 := range e.MappingsByArg0 {
		out = append(out, arg0)
	}
	sort.Strings(out)

	e.executables = out
	return out
}

func parseEnvArgs(input string) ([]string, []string, map[string]string, error) {
	const SHEBANG = "#!"

	var env, args []string
	var files map[string]string
	var err error

	// multiline script (with shebang)
	if newLineIndex := strings.IndexRune(input, '\n'); newLineIndex > -1 {
		if !strings.HasPrefix(input, SHEBANG) {
			return nil, nil, nil, errors.New("invalid shebang")
		}

		filename := filepath.Join(os.TempDir(), fmt.Sprintf("pimp-%d", time.Now().UTC().UnixNano()))

		s, ph := doPlaceholders(input[len(SHEBANG):newLineIndex])
		args, err = shellwords.Parse(s)
		if err != nil {
			return nil, nil, nil, err
		}
		args = append(undoPlaceholders(args, ph), filename)

		files = map[string]string{}
		files[filename] = input
	} else {
		s, ph := doPlaceholders(input)
		env, args, err = shellwords.ParseWithEnvs(s)
		if err != nil {
			return nil, nil, nil, err
		}
		args = undoPlaceholders(args, ph)
	}

	return env, args, files, nil
}

// shellwords.Parse and shellwords.ParseWithEnvs do not play nicely with Go
// templates {{ and }}. So we apply a little bit of dark magic to make
// everything work as expected. The idea is to replace {{ ... }} with
// ___pimp_X___ placeholders, parse the line, then to replace back.

var templateRegexp = regexp.MustCompile(`{{[^}]+}}`)
var placeholderRegexp = regexp.MustCompile(`___pimp_[0-9]+___`)

func doPlaceholders(input string) (string, map[string]string) {
	templatesByPlaceholder := map[string]string{}

	out := templateRegexp.ReplaceAllStringFunc(input, func(template string) string {
		placeholder := fmt.Sprintf("___pimp_%d___", len(templatesByPlaceholder))
		templatesByPlaceholder[placeholder] = template
		return placeholder
	})

	return out, templatesByPlaceholder
}

func undoPlaceholders(args []string, templatesByPlaceholder map[string]string) []string {
	out := make([]string, 0, len(args))

	for _, arg := range args {
		out = append(out, placeholderRegexp.ReplaceAllStringFunc(arg, func(placeholder string) string {
			return templatesByPlaceholder[placeholder]
		}))
	}

	return out
}
