package marshal

import "fmt"

func Go(input interface{}) (string, error) {
	return fmt.Sprintf("%#v", input), nil
}
