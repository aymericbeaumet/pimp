package kubernetes

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func Namespaces() ([]string, error) {
	k, err := newK8s()
	if err != nil {
		return nil, err
	}

	namespaces, err := k.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var out []string
	for _, ns := range namespaces.Items {
		out = append(out, ns.Name)
	}

	return out, nil
}

func newK8s() (*kubernetes.Clientset, error) {
	kubeconfigs := strings.Split(os.Getenv("KUBECONFIG"), ":")

	var kubeconfig string
	if len(kubeconfigs) > 0 {
		kubeconfig = kubeconfigs[0]
	}

	if len(kubeconfig) == 0 {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}
