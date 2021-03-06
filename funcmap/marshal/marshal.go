package marshal

import (
	"encoding/json"
	"encoding/xml"
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

		"YAML": func(input interface{}) (string, error) {
			out, err := yaml.Marshal(input)
			if err != nil {
				return "", err
			}
			return string(out), nil
		},
	}
}
