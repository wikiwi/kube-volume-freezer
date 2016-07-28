/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package kubernetes

import (
	"testing"

	api "k8s.io/kubernetes/pkg/api"
	errors "k8s.io/kubernetes/pkg/api/errors"
	"k8s.io/kubernetes/pkg/client/unversioned/testclient"
	"k8s.io/kubernetes/pkg/runtime"

	"github.com/wikiwi/kube-volume-freezer/pkg/util/diff"
	"github.com/wikiwi/kube-volume-freezer/pkg/util/reflect"
)

// NOTICE: test doesn't cover different namespace, because testclient doesn't
// distinquish different namespaces.

var fixture = []runtime.Object{
	&api.Pod{
		ObjectMeta: api.ObjectMeta{Namespace: "default", Name: "pod1", UID: "11111111-1111-1111-1111-111111111111", Labels: map[string]string{"abc": "123"}},
		Status:     api.PodStatus{PodIP: "127.0.0.1"},
		Spec:       api.PodSpec{NodeName: "node1"},
	},
	&api.Pod{
		ObjectMeta: api.ObjectMeta{Namespace: "default", Name: "pod2", UID: "22222222-2222-2222-2222-222222222222", Labels: map[string]string{"abc": "123"}},
		Status:     api.PodStatus{PodIP: "127.0.0.2"},
		Spec:       api.PodSpec{NodeName: "node2"},
	},
	&api.Pod{
		ObjectMeta: api.ObjectMeta{Namespace: "default", Name: "pod3", UID: "33333333-3333-3333-3333-333333333333", Labels: map[string]string{"def": "123"}},
		Status:     api.PodStatus{PodIP: "127.0.0.3"},
		Spec:       api.PodSpec{NodeName: "node3"},
	},
}

func TestDiscover(t *testing.T) {
	fakeClient := testclient.NewSimpleFake(fixture...)

	testScenarios := []struct {
		cfg        *DiscoveryConfig
		minionList MinionList
		err        error
	}{
		{
			cfg: &DiscoveryConfig{
				Namespace: "default",
				Selector:  "abc=123",
				Scheme:    "http",
				Port:      80,
			},
			minionList: MinionList{
				{Address: "http://127.0.0.1:80", NodeName: "node1"},
				{Address: "http://127.0.0.2:80", NodeName: "node2"},
			},
		},
		{
			cfg: &DiscoveryConfig{
				Namespace: "default",
				Selector:  "def=123",
				Scheme:    "http",
				Port:      80,
			},
			minionList: MinionList{
				{Address: "http://127.0.0.3:80", NodeName: "node3"},
			},
		},
		{
			cfg: &DiscoveryConfig{
				Namespace: "default",
				Selector:  "not=existing",
				Scheme:    "http",
				Port:      80,
			},
		},
	}
	for _, x := range testScenarios {
		svc, err := NewService(fakeClient, x.cfg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		minionList, err := svc.Discover()
		if x.err == nil && err != nil || !reflect.DeepDerivative(x.err, err) {
			t.Fatalf("unexpected error: %v", diff.ObjectDiffDerivative(x.err, err))
		}
		if !reflect.DeepEqual(x.minionList, minionList) {
			t.Fatalf("unexpected minionList: %v", diff.ObjectDiff(x.minionList, minionList))
		}
	}
}

func TestGetPodUID(t *testing.T) {
	fakeClient := testclient.NewSimpleFake(fixture...)
	fakeClient.PrependReactor("get", "*", func(action testclient.Action) (handled bool, ret runtime.Object, err error) {
		a := action.(testclient.GetAction)
		if a.GetName() == "not-existing" {
			return true, nil, errors.NewNotFound(api.Resource("pods"), a.GetName())
		}
		return false, nil, nil
	})

	testScenarios := []struct {
		namespace string
		name      string
		nfo       *PodInfo
		err       error
	}{
		{
			namespace: "default",
			name:      "pod1",
			nfo:       &PodInfo{UID: "11111111-1111-1111-1111-111111111111", NodeName: "node1"},
		},
		{
			namespace: "default",
			name:      "not-existing",
			err:       errors.NewNotFound(api.Resource("pods"), "not-existing"),
		},
	}
	for _, x := range testScenarios {
		svc, err := NewService(fakeClient, &DiscoveryConfig{})
		if err != nil {
			t.Fatalf("namespace: %q, name: %q, unexpected error: %v", x.namespace, x.name, err)
		}
		nfo, err := svc.GetPodInfo(x.namespace, x.name)
		if x.err == nil && err != nil || !reflect.DeepDerivative(x.err, err) {
			t.Fatalf("namespace: %q, name: %q, unexpected error: %v", x.namespace, x.name, diff.ObjectDiffDerivative(x.err, err))
		}
		if !reflect.DeepEqual(x.nfo, nfo) {
			t.Fatalf("namespace: %q, name: %q, unexpected PodInfo: %v", x.namespace, x.name, diff.ObjectDiff(x.nfo, nfo))
		}
	}
}
