// Package prelude contains the most commonly used utility functions
package prelude

import (
	"text/template"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"Exec":          Exec,
		"FZF":           FZF,
		"Head":          Head,
		"Reverse":       Reverse,
		"Sort":          Sort,
		"Tail":          Tail,
		"ToString":      ToString,
		"ToStringSlice": ToStringSlice,
		"Type":          Type,
	}
}
