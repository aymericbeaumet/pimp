package git

import "sort"

func GitRemotes() ([]string, error) {
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
}
