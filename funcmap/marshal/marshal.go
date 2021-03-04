package marshal

import (
	"encoding/json"
	"encoding/xml"

	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
)

func JSON(input interface{}) (string, error) {
	out, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func TOML(input interface{}) (string, error) {
	out, err := toml.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func XML(input interface{}) (string, error) {
	out, err := xml.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func YAML(input interface{}) (string, error) {
	out, err := yaml.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
