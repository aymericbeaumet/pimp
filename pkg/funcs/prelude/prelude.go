// Package prelude contains the most commonly used utility functions
package prelude

import (
	"text/template"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
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
		"type":          Type,
	}
}
