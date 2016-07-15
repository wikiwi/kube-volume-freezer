package rest

// Resource is the interface of a REST Resource.
type Resource interface {
	Register(s *Server)
}
