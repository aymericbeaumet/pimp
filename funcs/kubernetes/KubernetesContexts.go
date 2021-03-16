package kubernetes

func KubernetesContexts() ([]string, error) {
	_, config, err := createClientAndConfig()
	if err != nil {
		return nil, err
	}

	var out []string
	for name := range config.Contexts {
		out = append(out, name)
	}

	return out, nil
}
