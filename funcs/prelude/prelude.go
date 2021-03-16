// Package prelude contains the most commonly used utility functions
package prelude

import (
	"text/template"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"FZF":           FZF,
		"Head":          Head,
		"Tail":          Tail,
		"ToString":      ToString,
		"ToStringSlice": ToStringSlice,
	}
}
