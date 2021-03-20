package git

func References() ([]*Reference, error) {
	repo, err := Open()
	if err != nil {
		return nil, err
	}
	return repo.References()
}
