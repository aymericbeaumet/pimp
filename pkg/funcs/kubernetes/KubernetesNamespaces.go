package kubernetes

func KubernetesNamespaces() ([]*Namespace, error) {
	k, err := KubernetesOpen()
	if err != nil {
		return nil, err
	}
	return k.Namespaces()
}
