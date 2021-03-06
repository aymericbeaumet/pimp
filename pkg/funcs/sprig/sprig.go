// Package sprig wraps the Sprig template functions (https://masterminds.github.io/sprig/)
package sprig

import (
	"text/template"

	"github.com/Masterminds/sprig"
)

func FuncMap() template.FuncMap {
	return sprig.TxtFuncMap()
}
