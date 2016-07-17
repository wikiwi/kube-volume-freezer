/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package controllers

import (
	"testing"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/client/clienttest"
	"github.com/wikiwi/kube-volume-freezer/pkg/client/generic"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/fs"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/fs/fstest"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/volumes"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest/resttest"
)

func TestAuthFilter(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(rest.ForbiddenFilter, volumes.NewManager(fs.New())),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	podUID := "11111111-1111-1111-1111-zzzzzzzzzzzz"
	req := client.NewRequestOrDie("GET", "volumes/"+podUID, nil)
	exp := clienttest.ResponseExpectation{Code: 403}
	exp.DoAndValidateOrDie(t, client, req)
}

func TestGetVolumeListNotFound(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(nil, volumes.NewManager(fs.New())),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	podUID := "11111111-1111-1111-1111-zzzzzzzzzzzz"
	req := client.NewRequestOrDie("GET", "volumes/"+podUID, nil)
	exp := clienttest.ResponseExpectation{Code: 404}
	exp.DoAndValidateOrDie(t, client, req)
}

func TestGetVolumeListValidation(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(nil, volumes.NewManager(fstest.NewFake([]string{}))),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	req := client.NewRequestOrDie("GET", "volumes/invalid-uid", nil)
	exp := clienttest.ResponseExpectation{Code: 422}
	exp.DoAndValidateOrDie(t, client, req)
}

func TestGetVolumeList(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(nil, volumes.NewManager(fstest.NewFake([]string{
			volumes.PodsBasePath + "/11111111-1111-1111-1111-111111111111/volumes/kubernetes.io~empty-dir/empty1",
			volumes.PodsBasePath + "/11111111-1111-1111-1111-111111111111/volumes/kubernetes.io~empty-dir/empty2",
			volumes.PodsBasePath + "/11111111-1111-1111-1111-111111111111/volumes/kubernetes.io~gce-pd/pd1",
			volumes.PodsBasePath + "/11111111-1111-1111-1111-111111111111/volumes/kubernetes.io~nfs",
		}))),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	podUID := "11111111-1111-1111-1111-111111111111"
	req := client.NewRequestOrDie("GET", "volumes/"+podUID, nil)
	exp := clienttest.ResponseExpectation{
		Code: 200,
		Entity: &api.VolumeList{
			PodUID: podUID,
			Items:  []string{"empty1", "empty2", "pd1"},
		},
	}
	exp.DoAndValidateOrDie(t, client, req)
}

func TestGetVolume(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(nil, volumes.NewManager(fstest.NewFake([]string{
			volumes.PodsBasePath + "/11111111-1111-1111-1111-111111111111/volumes/kubernetes.io~empty-dir/empty1",
		}))),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	podUID := "11111111-1111-1111-1111-111111111111"
	name := "empty1"
	req := client.NewRequestOrDie("GET", "volumes/"+podUID+"/"+name, nil)
	exp := clienttest.ResponseExpectation{
		Code: 200,
		Entity: &api.Volume{
			PodUID: podUID,
			Name:   name,
		},
	}
	exp.DoAndValidateOrDie(t, client, req)
}

func TestGetVolumeNotFound(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(nil, volumes.NewManager(fstest.NewFake([]string{
			volumes.PodsBasePath + "/11111111-1111-1111-1111-111111111111/volumes/kubernetes.io~empty-dir/empty1",
		}))),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	testScenarios := []struct {
		podUID string
		name   string
	}{
		{podUID: "11111111-1111-1111-1111-zzzzzzzzzzzz", name: "empty1"},
		{podUID: "11111111-1111-1111-1111-111111111111", name: "noteempty1xisting"},
	}
	for _, x := range testScenarios {
		req := client.NewRequestOrDie("GET", "volumes/"+x.podUID+"/"+x.name, nil)
		exp := clienttest.ResponseExpectation{Code: 404}
		err := exp.DoAndValidate(client, req)
		if err != nil {
			t.Errorf("podUID: %q, name: %q, %v", x.podUID, x.name, err)
		}
	}
}

func TestGetVolumeValidation(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(nil, volumes.NewManager(fstest.NewFake([]string{}))),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	testScenarios := []struct {
		podUID string
		name   string
	}{
		{podUID: "11111111-1111-1111-1111-111111111111", name: "$awdfaw"},
		{podUID: "11111111-1111-1111-1111-invalid", name: "valid"},
	}
	for _, x := range testScenarios {
		req := client.NewRequestOrDie("GET", "volumes/"+x.podUID+"/"+x.name, nil)
		exp := clienttest.ResponseExpectation{Code: 422}
		err := exp.DoAndValidate(client, req)
		if err != nil {
			t.Errorf("podUID: %q, name: %q, %v", x.podUID, x.name, err)
		}
	}
}

func TestFreezeThawVolumeValidation(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(nil, volumes.NewManager(fstest.NewFake([]string{}))),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	testScenarios := []struct {
		podUID string
		name   string
	}{
		{podUID: "11111111-1111-1111-1111-111111111111", name: "$awdfaw"},
		{podUID: "11111111-1111-1111-1111-invalid", name: "valid"},
	}
	for _, x := range testScenarios {
		req := client.NewRequestOrDie("POST", "volumes/"+x.podUID+"/"+x.name, nil)
		exp := clienttest.ResponseExpectation{Code: 422}
		err := exp.DoAndValidate(client, req)
		if err != nil {
			t.Errorf("podUID: %q, name: %q, %v", x.podUID, x.name, err)
		}
	}
}

func TestFreezeThawVolumeNotFound(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(nil, volumes.NewManager(fstest.NewFake([]string{
			volumes.PodsBasePath + "/11111111-1111-1111-1111-111111111111/volumes/kubernetes.io~empty-dir/empty1",
		}))),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	testScenarios := []struct {
		podUID string
		name   string
	}{
		{podUID: "11111111-1111-1111-1111-zzzzzzzzzzzz", name: "empty1"},
		{podUID: "11111111-1111-1111-1111-111111111111", name: "notexisting"},
	}
	for _, x := range testScenarios {
		entity := api.FreezeThawRequest{Action: "freeze"}
		req := client.NewRequestOrDie("POST", "volumes/"+x.podUID+"/"+x.name, &entity)
		exp := clienttest.ResponseExpectation{Code: 404}
		err := exp.DoAndValidate(client, req)
		if err != nil {
			t.Errorf("podUID: %q, name: %q, %v", x.podUID, x.name, err)
		}
	}
}

func TestFreezeThawValidation(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(nil, volumes.NewManager(fstest.NewFake([]string{}))),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	testScenarios := []struct {
		podUID string
		name   string
		action string
	}{
		{podUID: "11111111-1111-1111-1111-111111111111", name: "$awdfaw", action: "freeze"},
		{podUID: "11111111-1111-1111-1111-invalid", name: "valid", action: "freeze"},
		{podUID: "11111111-1111-1111-1111-111111111111", name: "empty1", action: "invalid"},
	}
	for _, x := range testScenarios {
		entity := api.FreezeThawRequest{Action: x.action}
		req := client.NewRequestOrDie("POST", "volumes/"+x.podUID+"/"+x.name, &entity)
		exp := clienttest.ResponseExpectation{Code: 422}
		err := exp.DoAndValidate(client, req)
		if err != nil {
			t.Errorf("podUID: %q, name: %q, action: %q, %v", x.podUID, x.name, x.action, err)
		}
	}
}

func TestFreezeThaw(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(nil, volumes.NewManager(fstest.NewFake([]string{
			volumes.PodsBasePath + "/11111111-1111-1111-1111-111111111111/volumes/kubernetes.io~empty-dir/empty1",
		}))),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	testScenarios := []struct {
		podUID string
		name   string
		action string
	}{
		{podUID: "11111111-1111-1111-1111-111111111111", name: "empty1", action: "freeze"},
		{podUID: "11111111-1111-1111-1111-111111111111", name: "empty1", action: "thaw"},
	}
	for _, x := range testScenarios {
		entity := api.FreezeThawRequest{Action: x.action}
		req := client.NewRequestOrDie("POST", "volumes/"+x.podUID+"/"+x.name, &entity)
		exp := clienttest.ResponseExpectation{Code: 200}
		err := exp.DoAndValidate(client, req)
		if err != nil {
			t.Errorf("podUID: %q, name: %q, action: %q, %v", x.podUID, x.name, x.action, err)
		}
	}
}
