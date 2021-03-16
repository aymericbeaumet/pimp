package prelude

import "errors"

func Tail(input interface{}) (string, error) {
	lines := ToStringSlice(input)
	if len(lines) == 0 {
		return "", errors.New("empty input")
	}
	return lines[len(lines)-1], nil
}
