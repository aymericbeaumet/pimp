// Package semver contains Semantic Versioning helper functions (https://semver.org/)
package semver

import (
	"regexp"
	"text/template"

	"github.com/blang/semver/v4"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"SemverLatest": Latest,
		"SemverMajor":  Major,
		"SemverMinor":  Minor,
		"SemverParse":  Parse,
		"SemverPatch":  Patch,
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
