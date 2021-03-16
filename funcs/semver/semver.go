// Package semver contains Semantic Versioning helper functions (https://semver.org/)
package semver

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"text/template"

	"github.com/blang/semver/v4"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"SemverMajor": SemverMajor,
		"SemverMinor": SemverMinor,
		"SemverParse": SemverParse,
		"SemverPatch": SemverPatch,
	}
}

type Version struct {
	prefix  string
	version semver.Version
}

func (v Version) String() string {
	return v.prefix + v.version.String()
}

var prefixRegexp = regexp.MustCompile("^[^0-9]+")

func NewVersion(s string) (Version, error) {
	var prefix string

	if n := prefixRegexp.FindStringIndex(s); n != nil {
		prefix = s[:n[1]]
		s = s[n[1]:]
	}

	v, err := semver.Parse(s)
	if err != nil {
		return Version{}, err
	}

	return Version{prefix: prefix, version: v}, nil
}

func getLatestVersion(input interface{}) (*Version, error) {
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
