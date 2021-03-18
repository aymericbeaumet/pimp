package prelude

import "os"

func CD(dir string) (interface{}, error) {
	return nil, os.Chdir(dir)
}
