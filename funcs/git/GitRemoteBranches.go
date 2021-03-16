package git

func GitRemoteBranches() ([]*Reference, error) {
	repo, err := GitOpen()
	if err != nil {
		return nil, err
	}
	return repo.RemoteBranches()
}
