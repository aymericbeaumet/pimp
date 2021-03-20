package git

func Tags() ([]*Tag, error) {
	repo, err := Open()
	if err != nil {
		return nil, err
	}
	return repo.Tags()
}
