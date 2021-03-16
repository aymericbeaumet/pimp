package marshal

import "fmt"

func MarshalGo(input interface{}) (string, error) {
	return fmt.Sprintf("%#v", input), nil
}
