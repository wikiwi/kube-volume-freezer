package resttest

import (
	"net/http/httptest"

	"github.com/emicklei/go-restful"

	"github.com/wikiwi/kube-volume-freezer/pkg/rest"
)

// RunTestServer starts a REST Server with given Resource.
func RunTestServer(res rest.Resource) *httptest.Server {
	s := rest.NewServer()
	res.Register(s)
	return httptest.NewServer(s.Handler())
}

// RunFilterTestServer starts a REST Server with given filter and routing "/" to routeFunc.
func RunFilterTestServer(filter restful.FilterFunction, routeFunc restful.RouteFunction) *httptest.Server {
	ws := new(restful.WebService)
	ws.Path("/").Produces(restful.MIME_JSON)
	ws.Route(ws.GET("").To(routeFunc))
	ws.Filter(filter)
	s := rest.NewServer()
	s.Register(ws)
	return httptest.NewServer(s.Handler())
}
