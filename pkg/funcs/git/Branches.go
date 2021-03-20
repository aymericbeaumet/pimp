package git

func Branches() ([]*Reference, error) {
	repo, err := Open()
	if err != nil {
		return nil, err
	}
	return repo.Branches()
}
