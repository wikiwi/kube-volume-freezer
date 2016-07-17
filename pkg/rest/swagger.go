/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package rest

import (
	"github.com/emicklei/go-restful/swagger"
)

// RegisterSwagger exposes swagger definitions at /apidocs.json.
func RegisterSwagger(s *Server) {
	config := swagger.Config{
		WebServices: s.container.RegisteredWebServices(), // you control what services are visible
		DisableCORS: true,
		ApiPath:     "/apidocs.json",
	}
	swagger.RegisterSwaggerService(config, s.container)
}
