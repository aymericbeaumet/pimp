package marshal

import "encoding/xml"

func XMLIndent(input interface{}) (string, error) {
	out, err := xml.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}
