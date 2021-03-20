package git

func Remotes() ([]*Remote, error) {
	repo, err := Open()
	if err != nil {
		return nil, err
	}
	return repo.Remotes()
}
