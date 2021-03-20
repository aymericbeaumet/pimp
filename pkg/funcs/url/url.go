// Package url wraps part of the Go "url" package (https://golang.org/pkg/url/)
package url

import (
	"net/url"
	"text/template"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"URLParse":      Parse,
		"URLParseQuery": ParseQuery,
	}
}

type URL struct {
	url *url.URL
}

func (u URL) String() string {
	return u.url.String()
}

type QueryString struct {
	qs url.Values
}

func (qs QueryString) String() string {
	return qs.qs.Encode()
}
