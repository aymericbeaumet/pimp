package prelude

import "os"

func Ls() ([]*File, error) {
	cwd, err := PWD()
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(cwd)
	if err != nil {
		return nil, err
	}

	var out []*File
	for _, f := range files {
		out = append(out, &File{
			filename: f.Name(),
		})
	}
	return out, nil
}
