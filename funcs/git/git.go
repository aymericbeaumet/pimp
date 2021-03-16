// Package git contains Git helper functions (https://git-scm.com/)
package git

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	perrors "github.com/aymericbeaumet/pimp/errors"
	"github.com/go-git/go-git/v5"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"GitBranches":       GitBranches,
		"GitLocalBranches":  GitLocalBranches,
		"GitReferences":     GitReferences,
		"GitRemoteBranches": GitRemoteBranches,
		"GitRemotes":        GitRemotes,
		"GitTags":           GitTags,
	}
}

func openGitRepo(segments ...string) (*git.Repository, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	segments = append([]string{wd}, segments...)
	repopath := filepath.Join(segments...)

	for len(repopath) > 0 {
		repo, err := git.PlainOpen(repopath)
		if err == nil {
			return repo, nil
		}
		repopath = strings.TrimRight(path.Dir(repopath), "/")
	}

	return nil, perrors.NewFatalError(128, "fatal: not a git repository (or any of the parent directories): .git")
}
