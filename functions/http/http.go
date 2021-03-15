// Package http contains all the HTTP related functions
package http

import (
	"net/http"
	"text/template"
)

var httpClient http.Client

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"HttpGet":     HttpGet,
		"QueryString": QueryString,
		"URL":         URL,
	}
}
