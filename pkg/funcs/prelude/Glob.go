package prelude

import (
	"path/filepath"
)

func Glob(pattern string) ([]*File, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	out := make([]*File, 0, len(matches))
	for _, m := range matches {
		out = append(out, &File{filename: m})
	}
	return out, nil
}
