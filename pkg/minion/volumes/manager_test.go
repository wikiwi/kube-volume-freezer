package volumes

import (
	"testing"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/fs/fstest"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/diff"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/reflect"
)

var uid1 = "11111111-1111-1111-1111-111111111111"
var uid2 = "11111111-1111-1111-1111-222222222222"

var fixture = []string{
	PodsBasePath + "/" + uid1 + "/volumes/kubernetes.io~empty/vol1",
	PodsBasePath + "/" + uid1 + "/volumes/kubernetes.io~nfs/vol2",
	PodsBasePath + "/" + uid2 + "/volumes/kubernetes.io~123/vol3",
}

func TestList(t *testing.T) {
	m := NewManager(fstest.NewFake(fixture))

	testScenarios := []struct {
		uid        string
		volumeList *api.VolumeList
		err        error
	}{
		{uid: uid1, volumeList: &api.VolumeList{PodUID: uid1, Items: []string{"vol1", "vol2"}}},
		{uid: uid2, volumeList: &api.VolumeList{PodUID: uid2, Items: []string{"vol3"}}},
		{uid: "non-existing", err: &api.Error{Code: 404}},
	}
	for _, x := range testScenarios {
		volumeList, err := m.List(x.uid)
		if x.err == nil && err != nil || !reflect.DeepDerivative(x.err, err) {
			t.Fatalf("uid: %q, unexpected error: %v", x.uid, diff.ObjectDiffDerivative(x.err, err))
		}
		if !reflect.DeepEqual(x.volumeList, volumeList) {
			t.Fatalf("uid: %q, unexpected volume: %v", x.uid, diff.ObjectDiff(x.volumeList, volumeList))
		}
	}
}

func TestGet(t *testing.T) {
	m := NewManager(fstest.NewFake(fixture))

	testScenarios := []struct {
		uid     string
		volName string
		volume  *api.Volume
		err     error
	}{
		{uid: uid1, volName: "vol1", volume: &api.Volume{PodUID: uid1, Name: "vol1"}},
		{uid: uid1, volName: "non-existing", err: &api.Error{Code: 404}},
		{uid: "non-existing", volName: "vol1", err: &api.Error{Code: 404}},
	}
	for _, x := range testScenarios {
		volume, err := m.Get(x.uid, x.volName)
		if x.err == nil && err != nil || !reflect.DeepDerivative(x.err, err) {
			t.Fatalf("uid: %q, volName: %q, unexpected error: %v", x.uid, x.volName, diff.ObjectDiffDerivative(x.err, err))
		}
		if !reflect.DeepEqual(x.volume, volume) {
			t.Fatalf("uid: %q, volName: %q, unexpected volume: %v", x.uid, x.volName, diff.ObjectDiff(x.volume, volume))
		}
	}
}

func TestFreeze(t *testing.T) {
	testScenarios := []struct {
		uid     string
		volName string
		volume  *api.Volume
		err     error
	}{
		{uid: uid1, volName: "vol1", volume: &api.Volume{PodUID: uid1, Name: "vol1"}},
		{uid: uid1, volName: "non-existing", err: &api.Error{Code: 404}},
		{uid: "non-existing", volName: "vol1", err: &api.Error{Code: 404}},
	}
	for _, x := range testScenarios {
		fakeFS := fstest.NewFake(fixture)
		m := NewManager(fakeFS)
		volume, err := m.Freeze(x.uid, x.volName)
		if x.err == nil && err != nil || !reflect.DeepDerivative(x.err, err) {
			t.Fatalf("uid: %q, volName: %q, unexpected error: %v", x.uid, x.volName, diff.ObjectDiffDerivative(x.err, err))
		}
		if !reflect.DeepEqual(x.volume, volume) {
			t.Fatalf("uid: %q, volName: %q, unexpected volume: %v", x.uid, x.volName, diff.ObjectDiff(x.volume, volume))
		}
		if err == nil && len(fakeFS.Frozen) == 0 {
			t.Fatalf("uid: %q, volName: %q, volume was not frozen", x.uid, x.volName)
		}
	}
}

func TestThaw(t *testing.T) {
	testScenarios := []struct {
		uid     string
		volName string
		volume  *api.Volume
		err     error
	}{
		{uid: uid1, volName: "vol1", volume: &api.Volume{PodUID: uid1, Name: "vol1"}},
		{uid: uid1, volName: "non-existing", err: &api.Error{Code: 404}},
		{uid: "non-existing", volName: "vol1", err: &api.Error{Code: 404}},
	}
	for _, x := range testScenarios {
		fakeFS := fstest.NewFake(fixture)
		m := NewManager(fakeFS)
		volume, err := m.Thaw(x.uid, x.volName)
		if x.err == nil && err != nil || !reflect.DeepDerivative(x.err, err) {
			t.Fatalf("uid: %q, volName: %q, unexpected error: %v", x.uid, x.volName, diff.ObjectDiffDerivative(x.err, err))
		}
		if !reflect.DeepEqual(x.volume, volume) {
			t.Fatalf("uid: %q, volName: %q, unexpected volume: %v", x.uid, x.volName, diff.ObjectDiff(x.volume, volume))
		}
		if err == nil && len(fakeFS.Thawed) == 0 {
			t.Fatalf("uid: %q, volName: %q, volume was not thawed", x.uid, x.volName)
		}
	}
}
