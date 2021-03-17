package prelude

import "fmt"

func Printf(format string, args ...interface{}) interface{} {
	return fmt.Sprintf(format, args...)
}
