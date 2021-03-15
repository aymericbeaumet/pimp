package http

import "net/url"

func QueryString(input string) (url.Values, error) {
	return url.ParseQuery(input)
}
