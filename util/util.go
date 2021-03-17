package util

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func FilterEmptyStrings(input []string) []string {
	out := make([]string, 0, len(input))
	for _, i := range input {
		if trimmed := strings.TrimSpace(i); len(trimmed) > 0 {
			out = append(out, trimmed)
		}
	}
	return out
}

func NormalizePath(input string) (string, error) {
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

func StripShebang(input string) string {
	if !strings.HasPrefix(input, "#!") {
		return input
	}

	if newlineIndex := strings.IndexRune(input, '\n'); newlineIndex > -1 {
		return input[newlineIndex+1:]
	}

	return ""
}
