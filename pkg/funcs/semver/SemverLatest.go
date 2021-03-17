package semver

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/aymericbeaumet/pimp/pkg/funcs/git"
	"github.com/blang/semver/v4"
)

func SemverLatest(input interface{}) (*Version, error) {
	var versions []Version

	switch i := input.(type) {
	case string:
		p, err := NewVersion(i)
		if err != nil {
			return nil, err
		}
		versions = append(versions, p)

	case []string:
		for _, s := range i {
			p, err := NewVersion(s)
			if err != nil {
				return nil, err
			}
			versions = append(versions, p)
		}

	case Version:
		versions = append(versions, i)

	case []Version:
		versions = i

	case semver.Version:
		versions = append(versions, Version{version: i})

	case []semver.Version:
		for _, v := range i {
			versions = append(versions, Version{version: v})
		}

	case []*git.Tag:
		for _, t := range i {
			p, err := NewVersion(t.String())
			if err != nil {
				return nil, err
			}
			versions = append(versions, p)
		}

	default:
		return nil, fmt.Errorf("unsupported input type %s", reflect.TypeOf(i))
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].version.LT(versions[j].version)
	})

	if len(versions) == 0 {
		p, err := NewVersion("0.0.0")
		if err != nil {
			return nil, err
		}
		versions = append(versions, p)
	}

	return &versions[len(versions)-1], nil

}
