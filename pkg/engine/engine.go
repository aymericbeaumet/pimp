package engine

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/aymericbeaumet/pimp/pkg/util"
	"github.com/mattn/go-shellwords"
	"gopkg.in/yaml.v2"
)

type Pimpfile yaml.MapSlice

type Engine struct {
	Mappings map[string][]*Mapping `json:"mappings"`

	// used to cache calls to the Commands() method
	commands []string
}

type Mapping struct {
	CWD     string            `json:"cwd,omitempty"`
	Pattern []string          `json:"pattern"`
	Env     []string          `json:"env,omitempty"`
	Args    []string          `json:"args"`
	Files   map[string]string `json:"files,omitempty"`
}

func New() *Engine {
	return &Engine{
		Mappings: map[string][]*Mapping{},
	}
}

func (eng *Engine) LoadPimpfile(filename string, setCWD bool) error {
	normalized, err := util.NormalizePath(filename)
	if err != nil {
		return err
	}

	f, err := os.Open(normalized)
	if err != nil {
		return err
	}
	defer f.Close()

	var pimpfile Pimpfile
	if err := yaml.NewDecoder(f).Decode(&pimpfile); err != nil {
		return err
	}

	for _, item := range pimpfile {
		pattern, err := shellwords.Parse(item.Key.(string))
		if err != nil {
			return err
		}

		env, args, files, err := parseEnvArgs(item.Value)
		if err != nil {
			return err
		}

		var cwd string
		if setCWD {
			cwd = filepath.Dir(normalized)
		}

		eng.Mappings[pattern[0]] = append(eng.Mappings[pattern[0]], &Mapping{
			CWD:     cwd,
			Pattern: pattern,
			Env:     env,
			Args:    args,
			Files:   files,
		})
	}

	return nil
}

func (eng *Engine) Map(env []string, args []string) ([]string, []string, map[string]string, string) {
	if len(args) == 0 {
		return env, args, nil, ""
	}

	mappings, ok := eng.Mappings[args[0]]
	if !ok {
		return env, args, nil, ""
	}

	for _, mapping := range mappings {
		if mapping.Pattern[len(mapping.Pattern)-1] == "..." {
			pattern := mapping.Pattern[:len(mapping.Pattern)-1]
			lim := len(pattern)
			if lim > len(args) {
				lim = len(args)
			}
			if reflect.DeepEqual(pattern, args[:lim]) {
				return append(env[:], mapping.Env...), append(mapping.Args[:], args[lim:]...), mapping.Files, mapping.CWD
			}
		}

		if reflect.DeepEqual(mapping.Pattern, args) {
			return append(env[:], mapping.Env...), mapping.Args, mapping.Files, mapping.CWD
		}
	}

	return env, args, nil, ""
}

func (eng *Engine) JSON(w io.Writer) error {
	return json.NewEncoder(w).Encode(eng)
}

func (eng *Engine) Commands() []string {
	if eng.commands != nil {
		return eng.commands
	}

	out := make([]string, 0, len(eng.Mappings))
	for executable := range eng.Mappings {
		out = append(out, executable)
	}
	sort.Strings(out)

	eng.commands = out
	return out
}

func parseEnvArgs(input interface{}) ([]string, []string, map[string]string, error) {
	const SHEBANG = "#!"

	var env, args []string
	var files map[string]string
	var err error

	switch input := input.(type) {

	case []string:
		args = make([]string, 0, len(input))
		args = append(args, input...)

	case []interface{}:
		args = make([]string, 0, len(input))
		for _, i := range input {
			args = append(args, i.(string))
		}

	case string:
		// multiline script
		if newLineIndex := strings.IndexRune(input, '\n'); newLineIndex > -1 {
			filename := filepath.Join(os.TempDir(), fmt.Sprintf("pimp-%d", time.Now().UTC().UnixNano()))

			// with shebang
			if strings.HasPrefix(input, SHEBANG) {
				shebang, ph := doPlaceholders(input[len(SHEBANG):newLineIndex])
				shebangArgs, err := shellwords.Parse(shebang)
				if err != nil {
					return nil, nil, nil, err
				}
				args = append(undoPlaceholders(shebangArgs, ph), filename)
			} else {
				args = append(args, "pimp", "--run", filename)
			}

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

	default:
		return nil, nil, nil, fmt.Errorf("unsupported mapping value: %#v", input)

	}

	return env, args, files, nil
}

// shellwords.Parse and shellwords.ParseWithEnvs do not play nicely with Go
// templates {{ and }}. So we apply a little bit of dark magic to make
// everything work as expected. The idea is to replace {{ ... }} with
// ___pimp_X___ placeholders, parse the line, then to replace back.

var templateRegexp = regexp.MustCompile(`{{[^}]+}}`)
var placeholderRegexp = regexp.MustCompile("\x00pimp[0-9]+\x00")

func doPlaceholders(input string) (string, map[string]string) {
	templatesByPlaceholder := map[string]string{}

	out := templateRegexp.ReplaceAllStringFunc(input, func(template string) string {
		placeholder := fmt.Sprintf("\x00pimp%d\x00", len(templatesByPlaceholder))
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
