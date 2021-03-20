package git

func RemoteBranches() ([]*Reference, error) {
	repo, err := Open()
	if err != nil {
		return nil, err
	}
	return repo.RemoteBranches()
}
