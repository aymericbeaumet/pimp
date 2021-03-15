package marshal

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"text/template"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"JSON": func(input interface{}) (string, error) {
			out, err := json.Marshal(input)
			if err != nil {
				return "", err
			}
			return string(out), nil
		},

		"JSONIndent": func(input interface{}) (string, error) {
			out, err := json.MarshalIndent(input, "", "  ")
			if err != nil {
				return "", err
			}
			return string(out), nil
		},

		"TOML": func(input interface{}) (string, error) {
			out, err := toml.Marshal(input)
			if err != nil {
				return "", err
			}
			return string(out), nil
		},

		"XML": func(input interface{}) (string, error) {
			out, err := xml.Marshal(input)
			if err != nil {
				return "", err
			}
			return string(out), nil
		},

		"XMLIndent": func(input interface{}) (string, error) {
			out, err := xml.MarshalIndent(input, "", "  ")
			if err != nil {
				return "", err
			}
			return string(out), nil
		},

		"YAML": func(input interface{}) (string, error) {
			out, err := yaml.Marshal(input)
			if err != nil {
				return "", err
			}
			return string(out), nil
		},

		//

		"Go": func(input interface{}) (string, error) {
			return fmt.Sprintf("%#v", input), nil
		},

		"Shell": func(input interface{}) (string, error) {
			switch input := input.(type) {
			case string:
				return fmt.Sprintf("%#v", input), nil
			case []string:
				var out string
				for i, s := range input {
					if i > 0 {
						out += " "
					}
					out += fmt.Sprintf("%#v", s)
				}
				return out, nil
			default:
				return "", fmt.Errorf("unsupported type, received %#v", reflect.TypeOf(input))
			}
		},
	}
}
