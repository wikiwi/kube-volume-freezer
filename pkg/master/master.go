// Package master contains the implementation of the Master Server.
package master

import (
	"k8s.io/kubernetes/pkg/client/unversioned"

	"github.com/wikiwi/kube-volume-freezer/pkg/log"
	"github.com/wikiwi/kube-volume-freezer/pkg/master/controllers"
	"github.com/wikiwi/kube-volume-freezer/pkg/master/kubernetes"
	"github.com/wikiwi/kube-volume-freezer/pkg/master/volumes"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/client"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest"
	utilk8s "github.com/wikiwi/kube-volume-freezer/pkg/util/kubernetes"
)

// Options for starting the Master REST API Server.
type Options struct {
	// Token enables token-based authentication.
	Token string

	// MinionToken is the Token passed to the Minion API.
	MinionToken string

	// MinionSelector is a Kubernetes Pod Selector for finding Minion Pods.
	// E.g. app=kube-volume-freezer,component=minion
	MinionSelector string

	// MinionNamespace is the Kubernetes Namespace of the Minion Pods.
	MinionNamespace string

	// MinionPort is the port to the Minion API Server.
	MinionPort int

	// KubeClient is used for testing purposes.
	KubeClient unversioned.Interface
}

// NewRESTServer starts the Master REST API Server.
func NewRESTServer(opts *Options) (*rest.Server, error) {
	kubeClient := opts.KubeClient
	if kubeClient == nil {
		var err error
		kubeClient, err = utilk8s.NewDefaultClient()
		if err != nil {
			return nil, err
		}
	}
	k8s, err := kubernetes.NewService(
		kubeClient,
		&kubernetes.DiscoveryConfig{
			Selector:  opts.MinionSelector,
			Scheme:    "http",
			Namespace: opts.MinionNamespace,
			Port:      opts.MinionPort,
		})
	if err != nil {
		return nil, err
	}

	var authFilter = rest.NoOpFilter
	if len(opts.Token) > 0 {
		log.Instance().Info("Turn on authentication for clients")
		authFilter = rest.NewTokenAuthFilter(opts.Token)
	}

	if len(opts.MinionToken) > 0 {
		log.Instance().Info("Use token to authenticate to Minions")
	}

	server := rest.NewStandardServer()
	options := &client.Options{Token: opts.MinionToken}
	manager := volumes.NewManager(k8s, client.NewFactory(options))
	controllers.NewVolume(authFilter, manager).Register(server)
	rest.NewHealthzResource().Register(server)

	rest.RegisterSwagger(server)
	return server, nil
}
