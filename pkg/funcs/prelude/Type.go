package prelude

import "reflect"

func Type(input interface{}) reflect.Type {
	return reflect.TypeOf(input)
}
