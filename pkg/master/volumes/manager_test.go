package volumes

import (
	"testing"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/master/kubernetes"
	"github.com/wikiwi/kube-volume-freezer/pkg/master/kubernetes/k8stest"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/client"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/client/clienttest"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/diff"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/reflect"
)

func TestList(t *testing.T) {
	volumeList := &api.VolumeList{PodUID: "11111111-1111-1111-1111-111111111111", Items: []string{"v1", "v2"}}

	minionClient := clienttest.NewMock()
	minionClient.VolumesMock.On("List", volumeList.PodUID).Return(volumeList, nil)

	k8s := &k8stest.Fake{
		MinionList: kubernetes.MinionList{
			{NodeName: "node", Address: "http://test"},
		},
		PodInfoMap: map[string]*kubernetes.PodInfo{
			"default/pod":      {NodeName: "node", UID: volumeList.PodUID},
			"default/nominion": {NodeName: "nominion", UID: volumeList.PodUID},
		},
	}
	factory := &clienttest.FactoryFake{
		Clients: map[string]client.Interface{"http://test": minionClient},
	}
	manager := NewManager(k8s, factory.New)

	testScenarios := []struct {
		namespace, pod string
		volumeList     *api.VolumeList
		err            error
	}{
		{
			namespace: "default", pod: "pod",
			volumeList: volumeList,
		},
		{
			namespace: "default", pod: "notFound",
			err: &api.Error{Code: 404},
		},
		{
			namespace: "default", pod: "nominion",
			err: &api.Error{Code: 404},
		},
	}
	for _, x := range testScenarios {
		volumeList, err := manager.List(x.namespace, x.pod)
		if x.err == nil && err != nil || !reflect.DeepDerivative(x.err, err) {
			t.Fatalf("unexpected error: %v", diff.ObjectDiffDerivative(x.err, err))
		}
		if !reflect.DeepEqual(x.volumeList, volumeList) {
			t.Fatalf("unexpected volumeList: %v", diff.ObjectDiff(x.volumeList, volumeList))
		}
	}
}

func TestGet(t *testing.T) {
	volume := &api.Volume{PodUID: "11111111-1111-1111-1111-111111111111", Name: "volume"}

	minionClient := clienttest.NewMock()
	minionClient.VolumesMock.On("Get", volume.PodUID, "volume").Return(volume, nil)

	k8s := &k8stest.Fake{
		MinionList: kubernetes.MinionList{
			{NodeName: "node", Address: "http://test"},
		},
		PodInfoMap: map[string]*kubernetes.PodInfo{
			"default/pod":      {NodeName: "node", UID: volume.PodUID},
			"default/nominion": {NodeName: "nominion", UID: volume.PodUID},
		},
	}
	factory := &clienttest.FactoryFake{
		Clients: map[string]client.Interface{"http://test": minionClient},
	}
	manager := NewManager(k8s, factory.New)

	testScenarios := []struct {
		namespace, pod, volumeName string
		volume                     *api.Volume
		err                        error
	}{
		{
			namespace: "default", pod: "pod", volumeName: "volume",
			volume: volume,
		},
		{
			namespace: "default", pod: "notFound", volumeName: "volume",
			err: &api.Error{Code: 404},
		},
		{
			namespace: "default", pod: "nominion", volumeName: "volume",
			err: &api.Error{Code: 404},
		},
	}
	for _, x := range testScenarios {
		volume, err := manager.Get(x.namespace, x.pod, x.volumeName)
		if x.err == nil && err != nil || !reflect.DeepDerivative(x.err, err) {
			t.Fatalf("unexpected error: %v", diff.ObjectDiffDerivative(x.err, err))
		}
		if !reflect.DeepEqual(x.volume, volume) {
			t.Fatalf("unexpected volume: %v", diff.ObjectDiff(x.volume, volume))
		}
	}
}

func TestFreeze(t *testing.T) {
	volume := &api.Volume{PodUID: "11111111-1111-1111-1111-111111111111", Name: "volume"}

	minionClient := clienttest.NewMock()
	minionClient.VolumesMock.On("Freeze", volume.PodUID, "volume").Return(volume, nil)

	k8s := &k8stest.Fake{
		MinionList: kubernetes.MinionList{
			{NodeName: "node", Address: "http://test"},
		},
		PodInfoMap: map[string]*kubernetes.PodInfo{
			"default/pod":      {NodeName: "node", UID: volume.PodUID},
			"default/nominion": {NodeName: "nominion", UID: volume.PodUID},
		},
	}
	factory := &clienttest.FactoryFake{
		Clients: map[string]client.Interface{"http://test": minionClient},
	}
	manager := NewManager(k8s, factory.New)

	testScenarios := []struct {
		namespace, pod, volumeName string
		volume                     *api.Volume
		err                        error
	}{
		{
			namespace: "default", pod: "pod", volumeName: "volume",
			volume: volume,
		},
		{
			namespace: "default", pod: "notFound", volumeName: "volume",
			err: &api.Error{Code: 404},
		},
		{
			namespace: "default", pod: "nominion", volumeName: "volume",
			err: &api.Error{Code: 404},
		},
	}
	for _, x := range testScenarios {
		volume, err := manager.Freeze(x.namespace, x.pod, x.volumeName)
		if x.err == nil && err != nil || !reflect.DeepDerivative(x.err, err) {
			t.Fatalf("unexpected error: %v", diff.ObjectDiffDerivative(x.err, err))
		}
		if !reflect.DeepEqual(x.volume, volume) {
			t.Fatalf("unexpected volume: %v", diff.ObjectDiff(x.volume, volume))
		}
	}
}

func TestThaw(t *testing.T) {
	volume := &api.Volume{PodUID: "11111111-1111-1111-1111-111111111111", Name: "volume"}

	minionClient := clienttest.NewMock()
	minionClient.VolumesMock.On("Thaw", volume.PodUID, "volume").Return(volume, nil)

	k8s := &k8stest.Fake{
		MinionList: kubernetes.MinionList{
			{NodeName: "node", Address: "http://test"},
		},
		PodInfoMap: map[string]*kubernetes.PodInfo{
			"default/pod":      {NodeName: "node", UID: volume.PodUID},
			"default/nominion": {NodeName: "nominion", UID: volume.PodUID},
		},
	}
	factory := &clienttest.FactoryFake{
		Clients: map[string]client.Interface{"http://test": minionClient},
	}
	manager := NewManager(k8s, factory.New)

	testScenarios := []struct {
		namespace, pod, volumeName string
		volume                     *api.Volume
		err                        error
	}{
		{
			namespace: "default", pod: "pod", volumeName: "volume",
			volume: volume,
		},
		{
			namespace: "default", pod: "notFound", volumeName: "volume",
			err: &api.Error{Code: 404},
		},
		{
			namespace: "default", pod: "nominion", volumeName: "volume",
			err: &api.Error{Code: 404},
		},
	}
	for _, x := range testScenarios {
		volume, err := manager.Thaw(x.namespace, x.pod, x.volumeName)
		if x.err == nil && err != nil || !reflect.DeepDerivative(x.err, err) {
			t.Fatalf("unexpected error: %v", diff.ObjectDiffDerivative(x.err, err))
		}
		if !reflect.DeepEqual(x.volume, volume) {
			t.Fatalf("unexpected volume: %v", diff.ObjectDiff(x.volume, volume))
		}
	}
}
