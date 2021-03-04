package funcmap

import (
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/aymericbeaumet/pimp/funcmap/git"
	"github.com/aymericbeaumet/pimp/funcmap/kubernetes"
	"github.com/aymericbeaumet/pimp/funcmap/marshal"
	"github.com/aymericbeaumet/pimp/funcmap/misc"
)

func FuncMap() template.FuncMap {
	fm := template.FuncMap{
		// git
		"GitBranches":       git.GitBranches,
		"GitLocalBranches":  git.GitLocalBranches,
		"GitReferences":     git.GitReferences,
		"GitRemoteBranches": git.GitRemoteBranches,
		"GitRemotes":        git.GitRemotes,

		// kubernetes
		"KubernetesContexts":   kubernetes.KubernetesContexts,
		"KubernetesNamespaces": kubernetes.KubernetesNamespaces,

		// marshal
		"JSON": marshal.JSON,
		"TOML": marshal.TOML,
		"XML":  marshal.XML,
		"YAML": marshal.YAML,

		// miscellaneous
		"FZF":  misc.FZF,
		"Head": misc.Head,
		"Tail": misc.Tail,
	}

	// sprig
	for k, v := range sprig.TxtFuncMap() {
		if _, ok := fm[k]; ok {
			panic(fmt.Errorf("function `%s` from sprig already existy in the FuncMap", k))
		}
		fm[k] = v
	}

	return fm
}
