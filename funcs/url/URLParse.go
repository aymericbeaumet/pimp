package url

import "net/url"

func URLParse(input string) (*url.URL, error) {
	return url.Parse(input)
}
