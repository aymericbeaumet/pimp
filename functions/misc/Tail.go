package misc

import "errors"

func Tail(input []string) (string, error) {
	if len(input) == 0 {
		return "", errors.New("empty input")
	}
	return input[len(input)-1], nil
}
