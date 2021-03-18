package prelude

import "gopkg.in/yaml.v2"

func ToYAML(input interface{}) (string, error) {
	out, err := yaml.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
