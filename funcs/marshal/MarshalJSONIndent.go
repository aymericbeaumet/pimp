package marshal

import "encoding/json"

func MarshalJSONIndent(input interface{}) (string, error) {
	out, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}
