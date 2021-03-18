package prelude

import "encoding/xml"

func ToPrettyXML(input interface{}) (string, error) {
	out, err := xml.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", err
	}
	return string(out), nil
}
