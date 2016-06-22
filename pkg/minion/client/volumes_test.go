package client_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/diff"
)

func TestVolumeList(t *testing.T) {
	client, server, mux := setup()
	defer server.Close()

	mux.HandleFunc("/volumes/test", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"podUID": "test"}`))
	})
	volumeList, err := client.Volumes().List("test")
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

	mux.HandleFunc("/volumes/test/volname", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"name": "volname"}`))
	})
	volume, err := client.Volumes().Get("test", "volname")
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
		mux.HandleFunc("/volumes/test/volname", func(w http.ResponseWriter, r *http.Request) {
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
			volume, err = client.Volumes().Freeze("test", "volname")
		} else if x.action == "thaw" {
			volume, err = client.Volumes().Thaw("test", "volname")
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
