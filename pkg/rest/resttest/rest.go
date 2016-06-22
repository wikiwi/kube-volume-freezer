package resttest

import (
	"net/http/httptest"

	"github.com/emicklei/go-restful"

	"github.com/wikiwi/kube-volume-freezer/pkg/rest"
)

func RunTestServer(res rest.Resource) *httptest.Server {
	s := rest.NewServer()
	res.Register(s)
	return httptest.NewServer(s.Handler())
}

func RunFilterTestServer(filter restful.FilterFunction, route restful.RouteFunction) *httptest.Server {
	ws := new(restful.WebService)
	ws.Path("/").Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(route))
	ws.Filter(filter)
	s := rest.NewServer()
	s.Register(ws)
	return httptest.NewServer(s.Handler())
}
