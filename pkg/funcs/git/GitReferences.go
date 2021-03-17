package git

func GitReferences() ([]*Reference, error) {
	repo, err := GitOpen()
	if err != nil {
		return nil, err
	}
	return repo.References()
}
