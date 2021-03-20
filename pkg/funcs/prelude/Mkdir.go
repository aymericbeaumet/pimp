package prelude

import "os"

func Mkdir(path string) (*Directory, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, err
	}
	return &Directory{dirname: path}, nil
}
