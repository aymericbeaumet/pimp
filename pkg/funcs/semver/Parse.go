package semver

import "github.com/aymericbeaumet/pimp/pkg/funcs/prelude"

func Parse(input interface{}) (*Version, error) {
	v, err := NewVersion(prelude.ToString(input))
	if err != nil {
		return nil, err
	}

	return &v, nil
}
