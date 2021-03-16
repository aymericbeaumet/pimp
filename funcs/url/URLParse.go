package url

import "net/url"

func URLParse(input string) (*URL, error) {
	u, err := url.Parse(input)
	if err != nil {
		return nil, err
	}
	return &URL{url: u}, nil
}
