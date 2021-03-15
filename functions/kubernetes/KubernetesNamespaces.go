package kubernetes

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func KubernetesNamespaces() ([]string, error) {
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
}
