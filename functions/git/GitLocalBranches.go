package git

import (
	"sort"

	"github.com/go-git/go-git/v5/plumbing"
)

func GitLocalBranches() ([]string, error) {
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
}
