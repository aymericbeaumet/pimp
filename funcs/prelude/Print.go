package prelude

import (
	"fmt"
)

func Print(args ...interface{}) (string, error) {
	return fmt.Sprint(args...), nil
}
