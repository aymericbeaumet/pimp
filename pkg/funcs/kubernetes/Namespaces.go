package kubernetes

func Namespaces() ([]*Namespace, error) {
	k, err := Open()
	if err != nil {
		return nil, err
	}
	return k.Namespaces()
}
