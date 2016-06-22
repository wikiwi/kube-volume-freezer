package rest_test

import (
	"testing"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/client/clienttest"
	"github.com/wikiwi/kube-volume-freezer/pkg/client/generic"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest/resttest"
)

func TestHealthz(t *testing.T) {
	s := resttest.RunTestServer(&rest.HealthzResource{})
	defer s.Close()

	client := generic.NewOrDie(s.URL, nil)

	req := client.NewRequestOrDie("GET", "healthz", nil)
	exp := clienttest.ResponseExpectation{
		Code:   200,
		Entity: &api.Health{Status: "healthy"},
	}
	exp.DoAndValidateOrDie(t, client, req)
}
