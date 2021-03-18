package prelude

import "fmt"

func ToGo(input interface{}) (string, error) {
	return fmt.Sprintf("%#v", input), nil
}
