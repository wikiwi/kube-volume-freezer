package client

import (
	"net/http"

	"github.com/wikiwi/kube-volume-freezer/pkg/client/generic"
	"github.com/wikiwi/kube-volume-freezer/pkg/version"
)

var UserAgent = "kvf/" + version.Version

type Factory func(address string, httpClient *http.Client) (Interface, error)

type Interface interface {
	VolumesInterface
}

// Client to the master API.
type Client struct {
	*generic.Client

	// Services
	volumes VolumesService
}

func (c *Client) Volumes() VolumesService {
	return c.volumes
}

// NewOrDie is like New but panics upon error
func NewOrDie(address string, token string, httpClient *http.Client) *Client {
	c, err := New(address, token, httpClient)
	if err != nil {
		panic(err)
	}
	return c
}

//New returns a new Client.
func New(address string, token string, httpClient *http.Client) (*Client, error) {
	g, err := generic.New(address, httpClient)
	if err != nil {
		return nil, err
	}
	c := &Client{Client: g}
	c.Headers["User-Agent"] = UserAgent
	if token != "" {
		c.Headers["Authorization"] = "Bearer " + token
	}
	c.volumes = &volumesServiceImpl{client: c}
	return c, nil
}

func NewFactory(token string) Factory {
	return func(address string, httpClient *http.Client) (Interface, error) {
		c, err := New(address, token, httpClient)
		return c, err
	}
}
