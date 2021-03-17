package prelude

import (
	"reflect"
)

func Reverse(input interface{}) interface{} {
	value := reflect.ValueOf(input)

	if value.Kind() == reflect.Slice {
		out := reflect.MakeSlice(value.Type(), 0, value.Len())
		for i := value.Len() - 1; i >= 0; i-- {
			out = reflect.Append(out, value.Index(i))
		}
		return out.Interface()
	}

	return input
}
