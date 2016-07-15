// Package generic contains the shared Client code of the Minion and Master Client.
package generic

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
)

const (
	mediaType = "application/json"
)

// Client implements a generic basic client for consuming the projects REST API.
type Client struct {
	// Client is the native HTTP client all requests get delegated to.
	Client *http.Client

	// BaseURL is the base URL for relative paths.
	BaseURL *url.URL

	// Headers that are set on each request.
	Headers map[string]string
}

// NewRequest creates an API Request to urlStr which can be a relative string.
// If supplied the body will be included and encoded to JSON.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	return req, nil
}

// NewRequestOrDie is the same as NewRequest but panics on error.
func (c *Client) NewRequestOrDie(method, urlStr string, body interface{}) *http.Request {
	req, err := c.NewRequest(method, urlStr, body)
	if err != nil {
		panic(err)
	}
	return req
}

// Do sends an API request and returns the API response. If decodeTo is set, the response body will
// be decoded (expects JSON) into the variable. This method automatically parses an API error and
// returns the error.
func (c *Client) Do(req *http.Request, decodeTo interface{}) (*http.Response, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil {
			err = rerr
		}
	}()

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if decodeTo != nil {
		err := json.NewDecoder(resp.Body).Decode(decodeTo)
		if err != nil {
			return nil, err
		}
	}

	return resp, err
}

// CheckResponse checks the response for an API error and returns it.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	var apiError api.Error
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err := json.Unmarshal(data, &apiError)
		if err != nil {
			return errors.New(string(data))
		}
	}

	return &apiError
}

// New returns a new Client.
func New(address string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	parsed, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	c := &Client{Client: httpClient, BaseURL: parsed, Headers: map[string]string{}}
	return c, nil
}

// NewOrDie is the same as New but panics on error.
func NewOrDie(address string, httpClient *http.Client) *Client {
	c, err := New(address, httpClient)
	if err != nil {
		panic(err)
	}
	return c
}
