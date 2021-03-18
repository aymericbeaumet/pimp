package prelude

import "encoding/xml"

func ToXML(input interface{}) (string, error) {
	out, err := xml.Marshal(input)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
