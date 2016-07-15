package rest

import (
	"github.com/emicklei/go-restful"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/log"
)

// NewHealthzResource creates a new HealthzResource.
func NewHealthzResource() *HealthzResource {
	return &HealthzResource{}
}

// HealthzResource is a REST resource for reporting health status.
type HealthzResource struct{}

// Register adds Resource to the provided Server.
func (r *HealthzResource) Register(s *Server) {
	ws := new(restful.WebService)
	ws.Path("/healthz").
		Doc("show health status").
		Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(r.getHealth).
		Doc("show health status").
		Writes(api.Health{})) // on the response
	s.container.Add(ws)
}

func (r *HealthzResource) getHealth(request *restful.Request, response *restful.Response) {
	h := &api.Health{Status: "healthy"}
	err := response.WriteEntity(h)
	if err != nil {
		log.Instance().Error(err)
	}
}
