/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

// Package minion contains the implementation of the Minion Server.
package minion

import (
	"github.com/wikiwi/kube-volume-freezer/pkg/log"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/controllers"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/fs"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/volumes"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest"
)

// Options for starting the Minion REST API Server.
type Options struct {
	// Token enables token-based authentication.
	Token string

	// FS is used for testing purposes.
	FS fs.FileSystem
}

// NewRESTServer starts the Minion REST API Server.
func NewRESTServer(opts *Options) (*rest.Server, error) {
	server := rest.NewStandardServer()

	var authFilter = rest.NoOpFilter
	if len(opts.Token) > 0 {
		log.Instance().Info("Turn on authentication")
		authFilter = rest.NewTokenAuthFilter(opts.Token)
	}

	f := opts.FS
	if f == nil {
		f = fs.New()
	}

	manager := volumes.NewManager(f)
	controllers.NewVolume(authFilter, manager).Register(server)
	rest.NewHealthzResource().Register(server)
	rest.RegisterSwagger(server)

	return server, nil
}
