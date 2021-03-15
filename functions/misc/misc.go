// Package misc contains all the functions that don't fit in other packages
package misc

import (
	"text/template"
)

type ExecRet struct {
	Pid    int
	Status int
	Stdout string
	Stderr string
}

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"Exec": Exec,
		"FZF":  FZF,
		"Head": Head,
		"Tail": Tail,
	}
}
