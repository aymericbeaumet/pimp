// Package sprig wraps the Sprig template functions (https://masterminds.github.io/sprig/)
package sprig

import (
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

func FuncMap() template.FuncMap {
	in := sprig.TxtFuncMap()
	out := template.FuncMap{}

	for name, fn := range in {
		if _, ok := sprigIgnore[name]; ok {
			continue
		}

		out[name] = fn
	}

	return out
}
