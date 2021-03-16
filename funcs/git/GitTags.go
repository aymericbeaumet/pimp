package git

import (
	"sort"

	"github.com/go-git/go-git/v5/plumbing"
)

func GitTags() ([]string, error) {
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
}
