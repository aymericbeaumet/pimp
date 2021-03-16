package url

import (
	"net/url"
	"strings"
)

func URLParseQuery(input string) (*QueryString, error) {
	input = strings.TrimLeft(input, "?")
	qs, err := url.ParseQuery(input)
	if err != nil {
		return nil, err
	}
	return &QueryString{qs: qs}, nil
}
