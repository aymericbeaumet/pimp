package prelude

import "encoding/json"

func ToJSON(input interface{}) (string, error) {
	out, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
