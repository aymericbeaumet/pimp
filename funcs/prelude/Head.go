package prelude

import "errors"

func Head(input interface{}) (string, error) {
	lines := ToStringSlice(input)
	if len(lines) == 0 {
		return "", errors.New("empty input")
	}
	return lines[0], nil
}
