package git

func GitRemotes() ([]*Remote, error) {
	repo, err := GitOpen()
	if err != nil {
		return nil, err
	}
	return repo.Remotes()
}
