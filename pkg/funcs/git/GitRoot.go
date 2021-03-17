package git

func GitRoot() (string, error) {
	repo, err := GitOpen()
	if err != nil {
		return "", err
	}
	return repo.Root(), nil
}
