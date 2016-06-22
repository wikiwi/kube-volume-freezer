package controllers

import (
	"testing"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/client/clienttest"
	"github.com/wikiwi/kube-volume-freezer/pkg/client/generic"
	"github.com/wikiwi/kube-volume-freezer/pkg/master/volumes/volumestest"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest"
	"github.com/wikiwi/kube-volume-freezer/pkg/rest/resttest"
)

func TestAuthFilter(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(rest.ForbiddenFilter, &volumestest.ManagerMock{}),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	namespace := "default"
	podName := "test"
	req := client.NewRequestOrDie("GET", "volumes/"+namespace+"/"+podName, nil)
	exp := clienttest.ResponseExpectation{Code: 403}
	exp.DoAndValidateOrDie(t, client, req)
}

func TestGetVolumeList(t *testing.T) {
	namespace := "default"
	podName := "test"
	podUID := "11111111-1111-1111-1111-111111111111"
	volumeList := &api.VolumeList{
		PodUID: podUID,
		Items:  []string{"vol1", "vol2"},
	}

	m := &volumestest.ManagerMock{}
	m.On("List", namespace, podName).Return(volumeList, nil)
	s := resttest.RunTestServer(NewVolume(nil, m))
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	req := client.NewRequestOrDie("GET", "volumes/"+namespace+"/"+podName, nil)
	exp := clienttest.ResponseExpectation{Code: 200, Entity: volumeList}
	exp.DoAndValidateOrDie(t, client, req)
}

func TestGetVolumeListError(t *testing.T) {
	namespace := "default"
	podName := "test"
	err := &api.Error{Code: 404}

	m := &volumestest.ManagerMock{}
	m.On("List", namespace, podName).Return(nil, err)
	s := resttest.RunTestServer(NewVolume(nil, m))
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	req := client.NewRequestOrDie("GET", "volumes/"+namespace+"/"+podName, nil)
	exp := clienttest.ResponseExpectation{Code: 404, Entity: err}
	exp.DoAndValidateOrDie(t, client, req)
}

func TestGetVolumeListValidation(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(nil, &volumestest.ManagerMock{}),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	testScenarios := []struct {
		namespace string
		podName   string
	}{
		{namespace: "$nvalid", podName: "valid"},
		{namespace: "valid", podName: "$nvalid"},
	}
	for _, x := range testScenarios {
		req := client.NewRequestOrDie("GET", "volumes/"+x.namespace+"/"+x.podName, nil)
		exp := clienttest.ResponseExpectation{Code: 422}
		err := exp.DoAndValidate(client, req)
		if err != nil {
			t.Errorf("ns: %q, pod: %q, %v", x.namespace, x.podName, err)
		}
	}
}

func TestGetVolume(t *testing.T) {
	namespace := "default"
	podName := "test"
	volName := "name"
	podUID := "11111111-1111-1111-1111-111111111111"
	volume := &api.Volume{
		PodUID: podUID,
		Name:   volName,
	}

	m := &volumestest.ManagerMock{}
	m.On("Get", namespace, podName, volName).Return(volume, nil)
	s := resttest.RunTestServer(
		NewVolume(nil, m),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	req := client.NewRequestOrDie("GET", "volumes/"+namespace+"/"+podName+"/"+volName, nil)
	exp := clienttest.ResponseExpectation{Code: 200, Entity: volume}
	exp.DoAndValidateOrDie(t, client, req)
}

func TestGetVolumeError(t *testing.T) {
	namespace := "default"
	podName := "test"
	volName := "volume"
	err := &api.Error{Code: 404}

	m := &volumestest.ManagerMock{}
	m.On("Get", namespace, podName, volName).Return(nil, err)
	s := resttest.RunTestServer(NewVolume(nil, m))
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	req := client.NewRequestOrDie("GET", "volumes/"+namespace+"/"+podName+"/"+volName, nil)
	exp := clienttest.ResponseExpectation{Code: 404, Entity: err}
	exp.DoAndValidateOrDie(t, client, req)
}

func TestGetVolumeValidation(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(nil, &volumestest.ManagerMock{}),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	testScenarios := []struct {
		namespace string
		podName   string
		volName   string
	}{
		{namespace: "$nvalid", podName: "valid", volName: "valid"},
		{namespace: "valid", podName: "$nvalid", volName: "valid"},
		{namespace: "valid", podName: "valid", volName: "$nvalid"},
	}
	for _, x := range testScenarios {
		req := client.NewRequestOrDie("GET", "volumes/"+x.namespace+"/"+x.podName+"/"+x.volName, nil)
		exp := clienttest.ResponseExpectation{Code: 422}
		err := exp.DoAndValidate(client, req)
		if err != nil {
			t.Errorf("ns: %q, pod: %q, volume: %q, %v", x.namespace, x.podName, x.volName, err)
		}
	}
}

func TestFreezeThaw(t *testing.T) {
	namespace := "default"
	podName := "test"
	volName := "name"
	podUID := "11111111-1111-1111-1111-111111111111"

	testScenarios := []struct {
		action string
	}{
		{action: "freeze"},
		{action: "thaw"},
	}
	for _, x := range testScenarios {
		volume := &api.Volume{
			PodUID: podUID,
			Name:   volName,
		}

		m := &volumestest.ManagerMock{}
		if x.action == "freeze" {
			m.On("Freeze", namespace, podName, volName).Return(volume, nil)
		} else {
			m.On("Thaw", namespace, podName, volName).Return(volume, nil)
		}

		s := resttest.RunTestServer(
			NewVolume(nil, m),
		)
		defer s.Close()
		client := generic.NewOrDie(s.URL, nil)

		entity := api.FreezeThawRequest{Action: x.action}
		req := client.NewRequestOrDie("POST", "volumes/"+namespace+"/"+podName+"/"+volName, &entity)
		exp := clienttest.ResponseExpectation{Code: 200, Entity: volume}
		exp.DoAndValidateOrDie(t, client, req)
	}
}

func TestFreezeThawError(t *testing.T) {
	namespace := "default"
	podName := "test"
	volName := "volume"
	err := &api.Error{Code: 404}

	testScenarios := []struct {
		action string
	}{
		{action: "freeze"},
		{action: "thaw"},
	}
	for _, x := range testScenarios {
		m := &volumestest.ManagerMock{}
		if x.action == "freeze" {
			m.On("Freeze", namespace, podName, volName).Return(nil, err)
		} else {
			m.On("Thaw", namespace, podName, volName).Return(nil, err)
		}

		s := resttest.RunTestServer(
			NewVolume(nil, m),
		)
		defer s.Close()
		client := generic.NewOrDie(s.URL, nil)

		entity := api.FreezeThawRequest{Action: x.action}
		req := client.NewRequestOrDie("POST", "volumes/"+namespace+"/"+podName+"/"+volName, &entity)
		exp := clienttest.ResponseExpectation{Code: 404, Entity: err}
		exp.DoAndValidateOrDie(t, client, req)
	}
}

func TestFreezeThawValidation(t *testing.T) {
	s := resttest.RunTestServer(
		NewVolume(nil, &volumestest.ManagerMock{}),
	)
	defer s.Close()
	client := generic.NewOrDie(s.URL, nil)

	testScenarios := []struct {
		namespace string
		podName   string
		volName   string
	}{
		{namespace: "$nvalid", podName: "valid", volName: "valid"},
		{namespace: "valid", podName: "$nvalid", volName: "valid"},
		{namespace: "valid", podName: "valid", volName: "$nvalid"},
	}
	for _, x := range testScenarios {
		req := client.NewRequestOrDie("POST", "volumes/"+x.namespace+"/"+x.podName+"/"+x.volName, nil)
		exp := clienttest.ResponseExpectation{Code: 422}
		err := exp.DoAndValidate(client, req)
		if err != nil {
			t.Errorf("ns: %q, pod: %q, volume: %q, %v", x.namespace, x.podName, x.volName, err)
		}
	}
}
