// Package http wraps part of the Go "net/http" package (https://golang.org/pkg/net/http/)
package http

import (
	"net/http"
	"text/template"
)

var httpClient http.Client

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"HTTPGet": Get,
	}
}
