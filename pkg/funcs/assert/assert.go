package assert

import "text/template"

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"AssertDeep":  Deep,
		"AssertEqual": Equal,
	}
}
