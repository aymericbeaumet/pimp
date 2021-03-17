package prelude

import (
	"fmt"
)

func Print(args ...interface{}) string {
	return fmt.Sprint(args...)
}
