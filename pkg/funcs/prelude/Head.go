package prelude

import (
	"reflect"
)

func Head(input interface{}) interface{} {
	value := reflect.ValueOf(input)

	if value.Kind() == reflect.Slice {
		if value.Len() == 0 {
			return nil
		}
		return value.Index(0).Interface()
	}

	return input
}
