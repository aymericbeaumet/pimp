package git

func GitTags() ([]*Tag, error) {
	repo, err := GitOpen()
	if err != nil {
		return nil, err
	}
	return repo.Tags()
}
