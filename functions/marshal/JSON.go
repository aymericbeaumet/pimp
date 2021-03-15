package marshal

import "encoding/json"

func JSON(input interface{}) (string, error) {
	out, err := json.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
