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

type Options struct {
	Token           string
	MinionToken     string
	MinionSelector  string
	MinionNamespace string
	MinionPort      int

	// For testing purposes.
	KubeClient unversioned.Interface
}

func NewRestServer(opts *Options) (*rest.Server, error) {
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
	controllers.NewVolume(authFilter, volumes.NewManager(k8s, client.NewFactory(opts.MinionToken))).
		Register(server)
	rest.NewHealthzResource().Register(server)

	rest.RegisterSwagger(server)
	return server, nil
}
