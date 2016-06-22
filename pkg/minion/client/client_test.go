package client_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wikiwi/kube-volume-freezer/pkg/minion/client"
	"github.com/wikiwi/kube-volume-freezer/pkg/version"
)

func setup() (*client.Client, *httptest.Server, *http.ServeMux) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	c := client.NewOrDie(server.URL, "", nil)
	return c, server, mux
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("expected request method to be %q, but was %q", r.Method, expected)
	}
}

func TestUserAgent(t *testing.T) {
	c := client.NewOrDie("http://localhost", "", nil)
	r, err := c.NewRequest("GET", "relative", nil)
	if err != nil {
		t.Fatal(err)
	}
	expUA := "kvf/" + version.Version
	if r.UserAgent() != expUA {
		t.Fatalf("expected user-agent header to be %q but was %q", expUA, r.UserAgent())
	}
}

func TestAuthorizationHeader(t *testing.T) {
	c := client.NewOrDie("http://localhost", "token", nil)
	r, err := c.NewRequest("GET", "relative", nil)
	if err != nil {
		t.Fatal(err)
	}
	expAuth := "Bearer token"
	if auth := r.Header.Get("Authorization"); auth != expAuth {
		t.Fatalf("expected authorization header to be %q but was %q", expAuth, auth)
	}
}
