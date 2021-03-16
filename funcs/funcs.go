// Package funcs contains template functions which can be accessed from either
// Pimp scripts or Pimp template files.
//
// You can use this outside of the context of Pimp, to do so just pass the
// return value of this package's FuncMap function, or from any underlying
// package, to https://golang.org/pkg/text/template/#Template.Funcs.
package funcs

import (
	"fmt"
	"text/template"

	"github.com/aymericbeaumet/pimp/funcs/git"
	"github.com/aymericbeaumet/pimp/funcs/http"
	"github.com/aymericbeaumet/pimp/funcs/kubernetes"
	"github.com/aymericbeaumet/pimp/funcs/marshal"
	"github.com/aymericbeaumet/pimp/funcs/os"
	"github.com/aymericbeaumet/pimp/funcs/prelude"
	"github.com/aymericbeaumet/pimp/funcs/semver"
	"github.com/aymericbeaumet/pimp/funcs/sprig"
)

// FuncMap returns a merged map with all the functions supported by Pimp. Refer
// to the individual packages to read more about all the available functions.
func FuncMap() template.FuncMap {
	return merge(
		git.FuncMap(),
		http.FuncMap(),
		kubernetes.FuncMap(),
		marshal.FuncMap(),
		os.FuncMap(),
		prelude.FuncMap(),
		semver.FuncMap(),
		sprig.FuncMap(),
	)
}

func merge(fms ...template.FuncMap) template.FuncMap {
	out := template.FuncMap{}

	for _, fm := range fms {
		for k, v := range fm {
			if _, ok := out[k]; ok {
				panic(fmt.Errorf("implementation error: duplicate FuncMap function `%s`", k))
			}
			out[k] = v
		}
	}

	return out
}
