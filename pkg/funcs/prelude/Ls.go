package prelude

import "os"

func Ls(path ...string) ([]*File, error) {
	var name string

	if len(path) == 0 {
		cwd, err := PWD()
		if err != nil {
			return nil, err
		}
		name = cwd
	} else {
		name = path[0]
	}

	files, err := os.ReadDir(name)
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
