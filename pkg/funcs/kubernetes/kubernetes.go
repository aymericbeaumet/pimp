// Package kubernetes contains Kubernetes helper functions (https://kubernetes.io/)
package kubernetes

import (
	"context"
	"text/template"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd/api"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"KubernetesContexts":   KubernetesContexts,
		"KubernetesOpen":       KubernetesOpen,
		"KubernetesNamespaces": KubernetesNamespaces,
	}
}

type Context struct {
	name    string
	context *api.Context
}

func (c Context) String() string {
	return c.name
}

type Namespace struct {
	namespace *v1.Namespace
}

func (n Namespace) String() string {
	return n.namespace.Name
}

type Client struct {
	client *kubernetes.Clientset
	config *api.Config
}

func (k Client) Contexts() ([]*Context, error) {
	out := make([]*Context, 0, len(k.config.Contexts))
	for name, c := range k.config.Contexts {
		out = append(out, &Context{name: name, context: c})
	}
	return out, nil
}

func (k Client) Namespaces() ([]*Namespace, error) {
	namespaces, err := k.client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	out := make([]*Namespace, 0, len(namespaces.Items))
	for _, ns := range namespaces.Items {
		out = append(out, &Namespace{namespace: &ns})
	}
	return out, nil
}
