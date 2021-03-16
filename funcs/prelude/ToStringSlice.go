package prelude

import (
	"fmt"
	"reflect"
	"strings"
)

func ToStringSlice(input interface{}) []string {
	value := reflect.ValueOf(input)

	if value.Kind() == reflect.Slice {
		out := make([]string, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			out = append(out, fmt.Sprint(value.Index(i)))
		}
		return out
	}

	return strings.Split(fmt.Sprint(input), "\n")
}
