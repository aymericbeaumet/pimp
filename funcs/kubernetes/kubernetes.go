// Package kubernetes contains Kubernetes helper functions (https://kubernetes.io/)
package kubernetes

import (
	"text/template"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"KubernetesContexts":   KubernetesContexts,
		"KubernetesNamespaces": KubernetesNamespaces,
	}
}

func createClientAndConfig() (*kubernetes.Clientset, *api.Config, error) {
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)

	config, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, nil, err
	}

	clientConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, nil, err
	}

	return client, &config, nil
}
