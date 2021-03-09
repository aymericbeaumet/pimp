package normalize

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func Path(input string) (string, error) {
	input = strings.TrimSpace(input)
	if len(input) == 0 {
		return input, nil
	}

	expanded, err := homedir.Expand(input)
	if err != nil {
		return "", err
	}
	if strings.HasPrefix(expanded, "/") {
		return expanded, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(wd, expanded), nil
}
