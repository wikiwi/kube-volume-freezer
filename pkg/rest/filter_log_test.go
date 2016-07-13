package rest_test

import (
	"testing"

	"github.com/emicklei/go-restful"

	"github.com/wikiwi/kube-volume-freezer/pkg/api/errors"
	"github.com/wikiwi/kube-volume-freezer/pkg/client/clienttest"
	"github.com/wikiwi/kube-volume-freezer/pkg/client/generic"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest/resttest"
)

func TestLogFilter(t *testing.T) {
	var notFound = errors.NotFound("test")
	okFunc := func(request *restful.Request, response *restful.Response) {
		response.AddHeader("Foo", "bar")
		response.AddHeader("Foo", "bar2")
		response.AddHeader("Bar", "foo")
		response.WriteHeader(404)
		response.WriteEntity(notFound)
	}

	s := resttest.RunFilterTestServer(rest.LogFilter, okFunc)
	defer s.Close()

	client := generic.NewOrDie(s.URL, nil)
	req := client.NewRequestOrDie("GET", "/", nil)
	exp := clienttest.ResponseExpectation{
		Code:    404,
		Entity:  notFound,
		Headers: map[string][]string{"Foo": {"bar", "bar2"}, "Bar": {"foo"}},
	}
	exp.DoAndValidateOrDie(t, client, req)
}
