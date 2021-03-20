package prelude

import "os"

func Touch(name string) (*File, error) {
	_, err := os.Create(name)
	if err != nil {
		return nil, err
	}
	return &File{filename: name}, nil
}
