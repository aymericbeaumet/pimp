// Package prelude contains the most commonly used utility functions
package prelude

import (
	"text/template"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"Exec":          Exec,
		"Exit":          Exit,
		"FZF":           FZF,
		"Head":          Head,
		"Print":         Print,
		"Printf":        Printf,
		"Println":       Println,
		"Reverse":       Reverse,
		"Sort":          Sort,
		"Tail":          Tail,
		"ToString":      ToString,
		"ToStringSlice": ToStringSlice,
		"Type":          Type,
	}
}
