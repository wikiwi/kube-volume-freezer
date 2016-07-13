package minion

import (
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/controllers"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/fs"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/volumes"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest"
)

type Options struct {
	Token string

	// For testing purposes.
	FS fs.FileSystem
}

func NewRestServer(opts *Options) (*rest.Server, error) {
	server := rest.NewStandardServer()

	var authFilter = rest.NoOpFilter
	if len(opts.Token) > 0 {
		authFilter = rest.NewTokenAuthFilter(opts.Token)
	}

	f := opts.FS
	if f == nil {
		f = fs.New()
	}

	controllers.NewVolume(authFilter, volumes.NewManager(f)).
		Register(server)
	rest.NewHealthzResource().Register(server)
	rest.RegisterSwagger(server)

	return server, nil
}
