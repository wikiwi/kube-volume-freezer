/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package rest

import (
	"github.com/emicklei/go-restful"
)

// BlankLineFilter adds a blank line at the end of each request.
func BlankLineFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	defer func() {
		resp.Write([]byte("\n"))
	}()
	chain.ProcessFilter(req, resp)
}
