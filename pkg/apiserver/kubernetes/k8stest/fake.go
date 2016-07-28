/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package k8stest

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"

	"github.com/wikiwi/kube-volume-freezer/pkg/apiserver/kubernetes"
)

var _ kubernetes.Service = new(Fake)

// Fake implements a Kubernetes Service Fake.
type Fake struct {
	// MinionList is the static list of Minions returned on Discover.
	MinionList kubernetes.MinionList

	// PodInfoMap is a static map of Pods returned on GetPodInfo.
	// Key of map is of format namespace/name.
	PodInfoMap map[string]*kubernetes.PodInfo
}

// Discover returns static list of Minions.
func (f *Fake) Discover() (kubernetes.MinionList, error) {
	return f.MinionList, nil
}

// GetPodInfo returns PodInfo from static map.
func (f *Fake) GetPodInfo(namespace string, name string) (*kubernetes.PodInfo, error) {
	key := namespace + "/" + name
	if nfo, found := f.PodInfoMap[key]; found {
		return nfo, nil
	}
	return nil, errors.NewNotFound(api.Resource("pods"), name)
}
