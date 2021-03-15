// Package semver contains all the Semantic Versioning related functions (https://semver.org/)
package semver

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"text/template"

	"github.com/blang/semver/v4"
)

var firstNumber = regexp.MustCompile("^[^0-9]+")

type Version struct {
	Prefix string
	semver.Version
}

func (v Version) String() string {
	return v.Prefix + v.Version.String()
}

func parse(s string) (Version, error) {
	var prefix string

	if n := firstNumber.FindStringIndex(s); n != nil {
		prefix = s[:n[1]]
		s = s[n[1]:]
	}

	v, err := semver.Parse(s)
	if err != nil {
		return Version{}, err
	}

	return Version{Prefix: prefix, Version: v}, nil
}

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"Major": Major,
		"Minor": Minor,
		"Patch": Patch,
	}
}

func getLatestVersion(input interface{}) (*Version, error) {
	var versions []Version

	switch i := input.(type) {
	case string:
		p, err := parse(i)
		if err != nil {
			return nil, err
		}
		versions = append(versions, p)

	case []string:
		for _, s := range i {
			p, err := parse(s)
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
		versions = append(versions, Version{Version: i})

	case []semver.Version:
		for _, v := range i {
			versions = append(versions, Version{Version: v})
		}

	default:
		return nil, fmt.Errorf("unsupported input type %s", reflect.TypeOf(i))
	}

	sort.Slice(versions, func(i, j int) bool {
		return versions[i].Version.LT(versions[j].Version)
	})

	if len(versions) == 0 {
		p, err := parse("0.0.0")
		if err != nil {
			return nil, err
		}
		versions = append(versions, p)
	}

	return &versions[len(versions)-1], nil
}
