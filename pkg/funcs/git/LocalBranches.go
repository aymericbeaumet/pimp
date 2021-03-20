package git

func LocalBranches() ([]*Reference, error) {
	repo, err := Open()
	if err != nil {
		return nil, err
	}
	return repo.LocalBranches()
}
