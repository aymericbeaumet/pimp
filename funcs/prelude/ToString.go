package prelude

import (
	"fmt"
	"reflect"
	"strings"
)

func ToString(input interface{}) string {
	if value := reflect.ValueOf(input); value.Kind() == reflect.Slice {
		var sb strings.Builder
		for i := 0; i < value.Len(); i++ {
			if i > 0 {
				sb.WriteRune('\n')
			}
			fmt.Fprint(&sb, value.Index(i))
		}
		return sb.String()
	}

	return fmt.Sprint(input)
}
