package misc

import "errors"

func Head(input []string) (string, error) {
	if len(input) == 0 {
		return "", errors.New("empty input")
	}
	return input[0], nil
}
