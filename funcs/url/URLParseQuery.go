package url

import "net/url"

func URLParseQuery(input string) (url.Values, error) {
	return url.ParseQuery(input)
}
