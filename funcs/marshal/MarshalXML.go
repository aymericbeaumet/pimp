package marshal

import "encoding/xml"

func MarshalXML(input interface{}) (string, error) {
	out, err := xml.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
