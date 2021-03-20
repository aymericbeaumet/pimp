package csv

import "text/template"

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"CSVParse": Parse,
	}
}
