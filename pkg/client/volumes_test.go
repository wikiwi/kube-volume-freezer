/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/diff"
)

func setup() (*Client, *httptest.Server, *http.ServeMux) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	c := NewOrDie(server.URL, nil)
	return c, server, mux
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if expected != r.Method {
		t.Errorf("expected request method to be %q, but was %q", r.Method, expected)
	}
}

func TestVolumeList(t *testing.T) {
	client, server, mux := setup()
	defer server.Close()

	mux.HandleFunc("/volumes/default/podA", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		if r.Method != "GET" {
			t.Errorf("expected request method to be %q, but was %q", "GET", r.Method)
		}
		w.Write([]byte(`{"podUID": "test"}`))
	})
	volumeList, err := client.Volumes().List("default", "podA")
	if err != nil {
		t.Fatal(err)
	}
	if volumeList.PodUID != "test" {
		t.Fatalf("expected podUID to be %q but was %q", "test", volumeList.PodUID)
	}
}

func TestVolume(t *testing.T) {
	client, server, mux := setup()
	defer server.Close()

	mux.HandleFunc("/volumes/default/test/volname", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"name": "volname"}`))
	})
	volume, err := client.Volumes().Get("default", "test", "volname")
	if err != nil {
		t.Fatal(err)
	}
	if volume.Name != "volname" {
		t.Fatalf("expected name to be %q but was %q", "volname", volume.Name)
	}
}

func TestFreezeThaw(t *testing.T) {
	testScenarios := []struct {
		action string
	}{
		{"freeze"},
		{"thaw"},
	}
	for _, x := range testScenarios {
		client, server, mux := setup()
		defer server.Close()
		expected := &api.FreezeThawRequest{
			Action: x.action,
		}
		mux.HandleFunc("/volumes/default/test/volname", func(w http.ResponseWriter, r *http.Request) {
			testMethod(t, r, "POST")
			request := new(api.FreezeThawRequest)
			err := json.NewDecoder(r.Body).Decode(request)
			if err != nil {
				t.Fatalf("unable to decode JSON: %v", err)
			}
			if !reflect.DeepEqual(expected, request) {
				t.Fatalf("unexpected request: %v", diff.ObjectDiff(expected, request))
			}

			w.Write([]byte(`{"name": "volname"}`))
		})

		var volume *api.Volume
		var err error

		if x.action == "freeze" {
			volume, err = client.Volumes().Freeze("default", "test", "volname")
		} else if x.action == "thaw" {
			volume, err = client.Volumes().Thaw("default", "test", "volname")
		} else {
			panic("unexpected action: " + x.action)
		}
		if err != nil {
			t.Error(err)
		} else if volume.Name != "volname" {
			t.Errorf("expected name to be %q but was %q", "volname", volume.Name)
		}
	}
}
