// Package funcs contains template functions which can be accessed from either
// Pimp scripts or Pimp template files.
//
// You can use this outside of the context of Pimp, to do so just pass the
// return value of this package's FuncMap function, or from any underlying
// package, to https://golang.org/pkg/text/template/#Template.Funcs.
package funcs

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/aymericbeaumet/pimp/pkg/funcs/assert"
	"github.com/aymericbeaumet/pimp/pkg/funcs/csv"
	"github.com/aymericbeaumet/pimp/pkg/funcs/emoji"
	"github.com/aymericbeaumet/pimp/pkg/funcs/git"
	"github.com/aymericbeaumet/pimp/pkg/funcs/http"
	"github.com/aymericbeaumet/pimp/pkg/funcs/kubernetes"
	"github.com/aymericbeaumet/pimp/pkg/funcs/markdown"
	"github.com/aymericbeaumet/pimp/pkg/funcs/prelude"
	"github.com/aymericbeaumet/pimp/pkg/funcs/semver"
	"github.com/aymericbeaumet/pimp/pkg/funcs/sql"
	"github.com/aymericbeaumet/pimp/pkg/funcs/url"
)

// FuncMap returns a merged map with all the functions supported by Pimp. Refer
// to the individual packages to read more about all the available functions.
func FuncMap() template.FuncMap {
	out := merge(
		assert.FuncMap(),
		csv.FuncMap(),
		emoji.FuncMap(),
		git.FuncMap(),
		http.FuncMap(),
		kubernetes.FuncMap(),
		markdown.FuncMap(),
		prelude.FuncMap(),
		semver.FuncMap(),
		sql.FuncMap(),
		url.FuncMap(),
	)

	funcs := make([]string, 0, len(out)+1)
	out["funcs"] = func() []string { return funcs }
	for funcName, fn := range out {
		name := "func " + funcName + strings.TrimPrefix(reflect.TypeOf(fn).String(), "func")
		funcs = append(funcs, name)
	}
	sort.Strings(funcs)

	return out
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
