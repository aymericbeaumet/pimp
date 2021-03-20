package assert

import (
	"errors"
	"reflect"
)

func Deep(a, b interface{}) (bool, error) {
	if !reflect.DeepEqual(a, b) {
		return false, errors.New("not equal")
	}
	return true, nil
}
