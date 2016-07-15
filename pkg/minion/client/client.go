package client

import (
	"net/http"

	"github.com/wikiwi/kube-volume-freezer/pkg/client/generic"
	"github.com/wikiwi/kube-volume-freezer/pkg/version"
)

var UserAgent = "kvf/" + version.Version

type Factory func(address string) (Interface, error)

type Interface interface {
	VolumesInterface
}

type Options struct {
	HTTPClient *http.Client
	Token      string
}

// Client to the minion API.
type Client struct {
	*generic.Client

	// Services
	volumes VolumesService
}

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

//New returns a new Client.
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

func NewFactory(options *Options) Factory {
	return func(address string) (Interface, error) {
		c, err := New(address, options)
		return c, err
	}
}
