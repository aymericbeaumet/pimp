package prelude

import "fmt"

func Printf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
