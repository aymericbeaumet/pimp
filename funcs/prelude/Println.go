package prelude

import (
	"fmt"
)

func Println(args ...interface{}) (string, error) {
	return fmt.Sprintln(args...), nil
}
