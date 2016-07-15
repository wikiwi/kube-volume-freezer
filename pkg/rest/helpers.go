package rest

import (
	"github.com/emicklei/go-restful"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/api/errors"
)

// ReadEntityOrBadRequest parses body of request into readTo and
// replies with a 400 error message on failure.
func ReadEntityOrBadRequest(readTo interface{}, request *restful.Request, response *restful.Response) bool {
	err := request.ReadEntity(readTo)
	if err != nil {
		er := errors.BadRequest("Unable to parse Entity")
		if err := response.WriteHeaderAndEntity(er.Code, er); err != nil {
			panic(err)
		}
		return false
	}
	return true
}

// WriteValidationError creates an UnprocessableEntity Error and writes to response.
func WriteValidationError(issueList api.IssueList, response *restful.Response) {
	er := errors.UnprocessableEntity("Unable to validate request")
	for _, issue := range issueList {
		er.Append(issue)
	}
	if err := response.WriteHeaderAndEntity(er.Code, er); err != nil {
		panic(err)
	}
}

// RespondOrDie is a convenience function for responding to a request.
// If given err is nil, this function writes the entity and code to the response.
// Otherwise err and its error code is written to the response.
// It panics on an unexpected error.
func RespondOrDie(code int, entity interface{}, err error, response *restful.Response) {
	if err != nil {
		if apiErr, ok := err.(*api.Error); ok {
			code = apiErr.Code
			entity = apiErr
		} else {
			code = 500
			entity = errors.Unexpected(err.Error())
		}
	}
	if err := response.WriteHeaderAndEntity(code, entity); err != nil {
		panic(err)
	}
}
