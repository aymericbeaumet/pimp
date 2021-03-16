package prelude

import "reflect"

func Type(input interface{}) string {
	return reflect.TypeOf(input).String()
}
