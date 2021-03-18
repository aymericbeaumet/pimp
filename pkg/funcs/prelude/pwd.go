package prelude

import "os"

func PWD() (string, error) {
	return os.Getwd()
}
