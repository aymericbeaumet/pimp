package prelude

import "reflect"

func TypeOf(input interface{}) reflect.Type {
	return reflect.TypeOf(input)
}
