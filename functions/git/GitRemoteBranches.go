package git

import (
	"sort"

	"github.com/go-git/go-git/v5/plumbing"
)

func GitRemoteBranches() ([]string, error) {
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
}
