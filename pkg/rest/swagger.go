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
