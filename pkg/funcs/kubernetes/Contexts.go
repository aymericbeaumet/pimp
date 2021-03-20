package kubernetes

func Contexts() ([]*Context, error) {
	k, err := Open()
	if err != nil {
		return nil, err
	}
	return k.Contexts()
}
