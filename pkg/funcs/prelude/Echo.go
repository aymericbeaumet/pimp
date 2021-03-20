package prelude

import (
	"fmt"
)

func Echo(args ...interface{}) string {
	return fmt.Sprintln(args...)
}
