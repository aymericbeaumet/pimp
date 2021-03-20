package git

func Root() (string, error) {
	repo, err := Open()
	if err != nil {
		return "", err
	}
	return repo.Root(), nil
}
