package rest

import (
	"github.com/emicklei/go-restful"

	"github.com/wikiwi/kube-volume-freezer/pkg/api/errors"
	"github.com/wikiwi/kube-volume-freezer/pkg/log"
)

type tokenAuthFilter struct {
	Token string
}

func (f *tokenAuthFilter) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	header := req.HeaderParameter("Authorization")
	if header != "Bearer "+f.Token {
		er := errors.Forbidden("Unauthorized")
		err := resp.WriteHeaderAndEntity(er.Code, er)
		if err != nil {
			log.Instance().WithField("error", err).Error(err)
		}
		return
	}
	chain.ProcessFilter(req, resp)
}

func NewTokenAuthFilter(token string) restful.FilterFunction {
	f := &tokenAuthFilter{token}
	return f.Filter
}

// NewTokenAuthFilter returns a filter implementing authorization
// based on a static token.

func ForbiddenFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	er := errors.Forbidden("Unauthorized")
	err := resp.WriteHeaderAndEntity(er.Code, er)
	if err != nil {
		log.Instance().WithField("error", err).Error(err)
	}
}
