package semver

import "github.com/aymericbeaumet/pimp/funcs/prelude"

func SemverParse(input interface{}) (*Version, error) {
	v, err := NewVersion(prelude.ToString(input))
	if err != nil {
		return nil, err
	}

	return &v, nil
}
