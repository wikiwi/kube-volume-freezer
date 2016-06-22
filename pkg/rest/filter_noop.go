package rest

import (
	"github.com/emicklei/go-restful"
)

// NoOpFilter is a simple No-Op filter for go-restful.
func NoOpFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	chain.ProcessFilter(req, resp)
}
