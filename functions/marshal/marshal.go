// Package marshal contains all the Marshaling related functions
package marshal

import (
	"text/template"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"Go":         Go,
		"JSON":       JSON,
		"JSONIndent": JSONIndent,
		"Shell":      Shell,
		"TOML":       TOML,
		"XML":        XML,
		"XMLIndent":  XMLIndent,
		"YAML":       YAML,
	}
}
