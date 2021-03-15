package http

import "net/url"

func URL(input string) (*url.URL, error) {
	return url.Parse(input)
}
