package prelude

import (
	"fmt"
)

func Println(args ...interface{}) string {
	return fmt.Sprintln(args...)
}
