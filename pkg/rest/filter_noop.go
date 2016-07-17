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

// NoOpFilter is a simple No-Op filter for go-restful.
func NoOpFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	chain.ProcessFilter(req, resp)
}
