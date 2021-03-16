package prelude

import (
	"reflect"
)

func Tail(input interface{}) interface{} {
	value := reflect.ValueOf(input)

	if value.Kind() == reflect.Slice {
		if value.Len() == 0 {
			return nil
		}
		return value.Index(value.Len() - 1).Interface()
	}

	return input
}
