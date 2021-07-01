package driver

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient create a client of kubernetes
func NewClient(kubeConfigPath string) (kubernetes.Interface, error) {
	var config *rest.Config
	var err error

	if kubeConfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to build config from in-cluster: %w", err)
		}
	}

	return kubernetes.NewForConfig(config)
}

// NewWatcher create a watcher of ingress
func NewWatcher(clientset kubernetes.Interface) (*cache.ListWatch, error) {
	ingressWatcher := cache.NewListWatchFromClient(
		clientset.NetworkingV1().RESTClient(),
		"ingresses",
		corev1.NamespaceAll,
		fields.Everything(),
	)

	return ingressWatcher, nil
}
