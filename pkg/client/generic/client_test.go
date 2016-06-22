package generic

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/diff"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/reflect"
)

func setup() (*Client, *httptest.Server, *http.ServeMux) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	c := NewOrDie(server.URL, nil)
	return c, server, mux
}

func TestHeaders(t *testing.T) {
	agentName := "testAgent"
	c := NewOrDie("http://localhost", nil)
	c.Headers["User-Agent"] = agentName
	r, err := c.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	if r.UserAgent() != agentName {
		t.Fatalf("expected header user-agent to be %q but was %q", agentName, r.UserAgent())
	}
}

func TestBaseURL(t *testing.T) {
	expHost := "remotehost.com"
	expScheme := "http"
	c := NewOrDie(expScheme+"://"+expHost, nil)
	r, err := c.NewRequest("GET", "relative", nil)
	if err != nil {
		t.Fatal(err)
	}

	if r.URL.Scheme != expScheme {
		t.Fatalf("expected scheme to be %q but was %q", expScheme, r.URL.Scheme)
	}
	if r.Host != expHost {
		t.Fatalf("expected host to be %q but was %q", expHost, r.Host)
	}
}

func TestAPIError(t *testing.T) {
	client, server, mux := setup()
	defer server.Close()

	expected := &api.Error{
		Code: 404,
	}
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"code": 404}`))
	})

	r, err := client.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Do(r, nil)
	if apiErr, ok := err.(*api.Error); !ok {
		t.Fatalf("expected apiError but was %T", err)
	} else {
		if !reflect.DeepDerivative(expected, apiErr) {
			t.Fatalf("unexpected request: %v", diff.ObjectDiffDerivative(expected, apiErr))
		}
	}
}
