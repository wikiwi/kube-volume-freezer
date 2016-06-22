package rest

type Resource interface {
	Register(s *Server)
}
