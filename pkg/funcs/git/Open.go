package git

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	perrors "github.com/aymericbeaumet/pimp/pkg/errors"
	"github.com/go-git/go-git/v5"
)

func Open(segments ...string) (*Repository, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	segments = append([]string{wd}, segments...)
	repopath := filepath.Join(segments...)

	for len(repopath) > 0 {
		repo, err := git.PlainOpen(repopath)
		if err == nil {
			return &Repository{path: repopath, repository: repo}, nil
		}
		repopath = strings.TrimRight(path.Dir(repopath), "/")
	}

	return nil, perrors.NewFatalError(128, "fatal: not a git repository (or any of the parent directories): .git")
}
