package prelude

import "os"

func Tee(filename, input string) (string, error) {
	if err := os.WriteFile(filename, []byte(input), 0644); err != nil {
		return "", err
	}
	return input, nil
}
