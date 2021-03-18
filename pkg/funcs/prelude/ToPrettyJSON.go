package prelude

import "encoding/json"

func ToPrettyJSON(input interface{}) (string, error) {
	out, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}
