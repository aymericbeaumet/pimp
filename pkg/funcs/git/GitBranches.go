package git

func GitBranches() ([]*Reference, error) {
	repo, err := GitOpen()
	if err != nil {
		return nil, err
	}
	return repo.Branches()
}
