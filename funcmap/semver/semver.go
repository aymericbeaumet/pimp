package semver

import (
	"fmt"
	"reflect"
	"text/template"

	"github.com/blang/semver/v4"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"Major": func(input interface{}) (*semver.Version, error) {
			v, err := getLatestVersion(input)
			if err != nil {
				return nil, err
			}
			if err := v.IncrementMajor(); err != nil {
				return nil, err
			}
			return v, nil
		},

		"Minor": func(input interface{}) (*semver.Version, error) {
			v, err := getLatestVersion(input)
			if err != nil {
				return nil, err
			}
			if err := v.IncrementMinor(); err != nil {
				return nil, err
			}
			return v, nil
		},

		"Patch": func(input interface{}) (*semver.Version, error) {
			v, err := getLatestVersion(input)
			if err != nil {
				return nil, err
			}
			if err := v.IncrementPatch(); err != nil {
				return nil, err
			}
			return v, nil
		},
	}
}

func getLatestVersion(input interface{}) (*semver.Version, error) {
	var versions []semver.Version

	switch i := input.(type) {
	case string:
		p, err := semver.ParseTolerant(i)
		if err != nil {
			return nil, err
		}
		versions = append(versions, p)

	case []string:
		for _, e := range i {
			p, err := semver.ParseTolerant(e)
			if err != nil {
				return nil, err
			}
			versions = append(versions, p)
		}

	case semver.Version:
		versions = append(versions, i)

	case []semver.Version:
		versions = i

	default:
		return nil, fmt.Errorf("unsupported input type %s", reflect.TypeOf(i))
	}

	semver.Sort(versions)

	if len(versions) == 0 {
		p, err := semver.ParseTolerant("0.0.0")
		if err != nil {
			return nil, err
		}
		versions = append(versions, p)
	}

	return &versions[len(versions)-1], nil
}
