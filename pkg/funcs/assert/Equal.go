package assert

import "errors"

func Equal(a, b string) (bool, error) {
	if a != b {
		return false, errors.New("not equal")
	}
	return true, nil
}
