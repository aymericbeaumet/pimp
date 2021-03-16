// Package url wraps part of the Go "url" package (https://golang.org/pkg/url/)
package url

import "text/template"

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"URLParse":      URLParse,
		"URLParseQuery": URLParseQuery,
	}
}
