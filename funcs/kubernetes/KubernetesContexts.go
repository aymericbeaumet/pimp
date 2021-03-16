package kubernetes

func KubernetesContexts() ([]*Context, error) {
	k, err := KubernetesOpen()
	if err != nil {
		return nil, err
	}
	return k.Contexts()
}
