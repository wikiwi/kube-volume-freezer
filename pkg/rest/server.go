package rest

import (
	"net/http"

	"github.com/emicklei/go-restful"

	"github.com/wikiwi/kube-volume-freezer/pkg/log"
)

// Server is a RESTful HTTP Server.
type Server struct {
	container *restful.Container
}

func NewStandardServer() *Server {
	container := restful.NewContainer()
	container.DoNotRecover(true)
	container.Filter(RecoverFilter)
	container.Filter(LogFilter)
	container.Filter(BlankLineFilter)
	container.Filter(container.OPTIONSFilter)
	container.EnableContentEncoding(true)
	return &Server{container: container}
}

func NewServer() *Server {
	container := restful.NewContainer()
	container.DoNotRecover(true)
	return &Server{container: container}
}

func (s *Server) Register(ws *restful.WebService) {
	s.container.Add(ws)
}

// Handler returns the http handler.
func (s *Server) Handler() http.Handler {
	return s.container
}

// ListenAndServe starts the HTTP Server and blocks.
func (s *Server) ListenAndServe(listen string) error {
	log.Instance().Printf("start listening on %s", listen)
	server := &http.Server{Addr: listen, Handler: s.container}
	return server.ListenAndServe()
}
