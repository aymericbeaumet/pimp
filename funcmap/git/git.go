package git

import (
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/aymericbeaumet/pimp/funcmap/errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"GitBranches": func() ([]string, error) {
			repo, err := openGitRepo()
			if err != nil {
				return nil, err
			}

			iter, err := repo.References()
			if err != nil {
				return nil, err
			}

			out := []string{}
			if err := iter.ForEach(func(reference *plumbing.Reference) error {
				name := reference.Name()
				if name.IsBranch() || name.IsRemote() {
					out = append(out, reference.Name().Short())
				}
				return nil
			}); err != nil {
				return nil, err
			}
			sort.Strings(out)

			return out, nil
		},

		"GitLocalBranches": func() ([]string, error) {
			repo, err := openGitRepo()
			if err != nil {
				return nil, err
			}

			iter, err := repo.Branches()
			if err != nil {
				return nil, err
			}

			out := []string{}
			if err := iter.ForEach(func(branch *plumbing.Reference) error {
				out = append(out, branch.Name().Short())
				return nil
			}); err != nil {
				return nil, err
			}
			sort.Strings(out)

			return out, nil
		},

		"GitReferences": func() ([]string, error) {
			repo, err := openGitRepo()
			if err != nil {
				return nil, err
			}

			iter, err := repo.References()
			if err != nil {
				return nil, err
			}

			out := []string{}
			if err := iter.ForEach(func(reference *plumbing.Reference) error {
				out = append(out, reference.Name().Short())
				return nil
			}); err != nil {
				return nil, err
			}
			sort.Strings(out)

			return out, nil
		},

		"GitRemoteBranches": func() ([]string, error) {
			repo, err := openGitRepo()
			if err != nil {
				return nil, err
			}

			iter, err := repo.References()
			if err != nil {
				return nil, err
			}

			out := []string{}
			if err := iter.ForEach(func(reference *plumbing.Reference) error {
				name := reference.Name()
				if name.IsRemote() {
					out = append(out, reference.Name().Short())
				}
				return nil
			}); err != nil {
				return nil, err
			}
			sort.Strings(out)

			return out, nil
		},

		"GitRemotes": func() ([]string, error) {
			repo, err := openGitRepo()
			if err != nil {
				return nil, err
			}

			remotes, err := repo.Remotes()
			if err != nil {
				return nil, err
			}

			out := []string{}
			for _, remote := range remotes {
				out = append(out, remote.Config().Name)
			}
			sort.Strings(out)

			return out, nil
		},

		"GitTags": func() ([]string, error) {
			repo, err := openGitRepo()
			if err != nil {
				return nil, err
			}

			iter, err := repo.Tags()
			if err != nil {
				return nil, err
			}

			out := []string{}
			if err := iter.ForEach(func(reference *plumbing.Reference) error {
				tag, err := repo.TagObject(reference.Hash())
				if err != nil {
					return err
				}
				out = append(out, tag.Name)
				return nil
			}); err != nil {
				return nil, err
			}
			sort.Strings(out)

			return out, nil
		},
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

	return nil, errors.NewFatalError(128, "fatal: not a git repository (or any of the parent directories): .git")
}
