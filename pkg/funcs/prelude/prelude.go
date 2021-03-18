// Package prelude contains the most commonly used utility functions
package prelude

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
	"typeOf":       {},
}

func FuncMap() template.FuncMap {
	out := template.FuncMap{
		"exec":          Exec,
		"exit":          Exit,
		"fzf":           FZF,
		"head":          Head,
		"print":         Print,
		"printf":        Printf,
		"println":       Println,
		"reverse":       Reverse,
		"sort":          Sort,
		"tail":          Tail,
		"toGo":          ToGo,
		"toJSON":        ToJSON,
		"toPrettyJSON":  ToPrettyJSON,
		"toPrettyXML":   ToPrettyXML,
		"toShell":       ToShell,
		"toString":      ToString,
		"toStringSlice": ToStringSlice,
		"toTOML":        ToTOML,
		"toXML":         ToXML,
		"toYAML":        ToYAML,
		"typeOf":        TypeOf,
	}

	for name, fn := range sprig.TxtFuncMap() {
		if _, ok := sprigIgnore[name]; ok {
			continue
		}
		out[name] = fn
	}

	return out
}
