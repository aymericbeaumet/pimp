package prelude

import "github.com/pelletier/go-toml"

func ToTOML(input interface{}) (string, error) {
	out, err := toml.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
