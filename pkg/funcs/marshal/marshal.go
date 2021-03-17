// Package marshal contains Marshaling related functions (JSON, YAML, etc)
package marshal

import (
	"text/template"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"MarshalGo":         MarshalGo,
		"MarshalJSON":       MarshalJSON,
		"MarshalJSONIndent": MarshalJSONIndent,
		"MarshalShell":      MarshalShell,
		"MarshalTOML":       MarshalTOML,
		"MarshalXML":        MarshalXML,
		"MarshalXMLIndent":  MarshalXMLIndent,
		"MarshalYAML":       MarshalYAML,
	}
}
