package kubernetes

import (
	"context"
	"text/template"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"KubernetesContexts": func() ([]string, error) {
			_, config, err := createClientAndConfig()
			if err != nil {
				return nil, err
			}

			var out []string
			for name := range config.Contexts {
				out = append(out, name)
			}

			return out, nil
		},

		"KubernetesNamespaces": func() ([]string, error) {
			client, _, err := createClientAndConfig()
			if err != nil {
				return nil, err
			}

			namespaces, err := client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				return nil, err
			}

			var out []string
			for _, ns := range namespaces.Items {
				out = append(out, ns.Name)
			}

			return out, nil
		},
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
