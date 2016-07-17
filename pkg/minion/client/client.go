/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

// Package client implements a client to the kube-volume-freezer Minion API.
package client

import (
	"net/http"

	"github.com/wikiwi/kube-volume-freezer/pkg/client/generic"
	"github.com/wikiwi/kube-volume-freezer/pkg/version"
)

// UserAgent that is sent with the HTTP Header on each request.
var UserAgent = "kvf/" + version.Version

// Factory creates instances of Client.
type Factory func(address string) (Interface, error)

// Interface of the Client.
type Interface interface {
	VolumesInterface
}

// Options is used for creating an instance of Client.
type Options struct {
	HTTPClient *http.Client
	Token      string
}

// Client performs requests to the Minion API.
type Client struct {
	*generic.Client

	// Services
	volumes VolumesService
}

// Volumes returns a Service to manipulate Volume Resources.
func (c *Client) Volumes() VolumesService {
	return c.volumes
}

// NewOrDie is like New but panics upon error
func NewOrDie(address string, opts *Options) *Client {
	c, err := New(address, opts)
	if err != nil {
		panic(err)
	}
	return c
}

// New returns a new Client.
func New(address string, opts *Options) (*Client, error) {
	if opts == nil {
		opts = new(Options)
	}
	g, err := generic.New(address, opts.HTTPClient)
	if err != nil {
		return nil, err
	}
	c := &Client{Client: g}
	c.Headers["User-Agent"] = UserAgent
	if opts.Token != "" {
		c.Headers["Authorization"] = "Bearer " + opts.Token
	}
	c.volumes = &volumesServiceImpl{client: c}
	return c, nil
}

// NewFactory returns a new Factory.
func NewFactory(options *Options) Factory {
	return func(address string) (Interface, error) {
		c, err := New(address, options)
		return c, err
	}
}
