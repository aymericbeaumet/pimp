package git

func GitLocalBranches() ([]*Reference, error) {
	repo, err := GitOpen()
	if err != nil {
		return nil, err
	}
	return repo.LocalBranches()
}
