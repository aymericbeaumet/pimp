// Package os wraps part of the Go "os" packages (https://golang.org/pkg/os/)
package os

import "text/template"

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"OSExec": OSExec,
	}
}
