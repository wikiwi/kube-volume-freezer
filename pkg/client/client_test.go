package client

import (
	"testing"

	"github.com/wikiwi/kube-volume-freezer/pkg/version"
)

func TestUserAgent(t *testing.T) {
	c := NewOrDie("http://localhost", "", nil)
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
	c := NewOrDie("http://localhost", "token", nil)
	r, err := c.NewRequest("GET", "relative", nil)
	if err != nil {
		t.Fatal(err)
	}
	expAuth := "Bearer token"
	if auth := r.Header.Get("Authorization"); auth != expAuth {
		t.Fatalf("expected authorization header to be %q but was %q", expAuth, auth)
	}
}
