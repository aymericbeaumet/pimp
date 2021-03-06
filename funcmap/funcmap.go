package funcmap

import (
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/aymericbeaumet/pimp/funcmap/git"
	"github.com/aymericbeaumet/pimp/funcmap/http"
	"github.com/aymericbeaumet/pimp/funcmap/kubernetes"
	"github.com/aymericbeaumet/pimp/funcmap/marshal"
	"github.com/aymericbeaumet/pimp/funcmap/misc"
)

func FuncMap() template.FuncMap {
	return merge(
		git.FuncMap(),
		http.FuncMap(),
		kubernetes.FuncMap(),
		marshal.FuncMap(),
		misc.FuncMap(),
		sprig.TxtFuncMap(),
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
