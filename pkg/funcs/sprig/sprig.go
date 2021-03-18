// Package sprig wraps the Sprig template functions (https://masterminds.github.io/sprig/)
package sprig

import (
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
)

var sprigIgnore = map[string]struct{}{
	"reverse":      {},
	"toJson":       {},
	"toPrettyJson": {},
	"toRawJson":    {},
	"toString":     {},
}

const sprigMustPrefix = "must"

func FuncMap() template.FuncMap {
	in := sprig.TxtFuncMap()
	out := template.FuncMap{}

	for name, fn := range in {
		// skip all the functions from the ignore list
		if _, ok := sprigIgnore[name]; ok {
			continue
		}

		// skip all the "must" functions
		if strings.HasPrefix(name, sprigMustPrefix) {
			continue
		}

		// try to find a must variant and use it if found
		mustVariant := sprigMustPrefix + strings.ToUpper(name[0:1]) + name[1:]
		if _, ok := in[mustVariant]; ok {
			name = mustVariant
		}

		out[name] = fn
	}

	return out
}
