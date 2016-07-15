// Package kubernetes contains tools for interfacing with kubernetes.
package kubernetes

import (
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/client/unversioned/clientcmd"
)

// NewDefaultClient returns a new default Kubernetes client.
func NewDefaultClient() (*unversioned.Client, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := unversioned.New(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}
