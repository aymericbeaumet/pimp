// Package semver contains all the Semantic Versioning related functions (https://semver.org/)
package semver

import (
	"fmt"
	"reflect"
	"text/template"

	"github.com/blang/semver/v4"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"Major": Major,
		"Minor": Minor,
		"Patch": Patch,
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
