package marshal

import (
	"fmt"
	"reflect"
)

func MarshalShell(input interface{}) (string, error) {
	switch input := input.(type) {
	case string:
		return fmt.Sprintf("%#v", input), nil
	case []string:
		var out string
		for i, s := range input {
			if i > 0 {
				out += " "
			}
			out += fmt.Sprintf("%#v", s)
		}
		return out, nil
	default:
		return "", fmt.Errorf("unsupported type, received %#v", reflect.TypeOf(input))
	}
}
