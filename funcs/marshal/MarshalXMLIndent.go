package marshal

import "encoding/xml"

func MarshalXMLIndent(input interface{}) (string, error) {
	out, err := xml.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}
