package rest

import (
	"github.com/emicklei/go-restful"
)

// BlankLineFilter adds a blank line at the end of the request.
func BlankLineFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	defer func() {
		resp.Write([]byte("\n"))
	}()
	chain.ProcessFilter(req, resp)
}
