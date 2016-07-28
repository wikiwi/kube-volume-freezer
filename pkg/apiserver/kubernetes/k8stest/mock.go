/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package k8stest

import (
	"github.com/stretchr/testify/mock"

	"github.com/wikiwi/kube-volume-freezer/pkg/apiserver/kubernetes"
)

var _ kubernetes.Service = new(Mock)

// Mock implements a Kubernetes Service Mock.
type Mock struct {
	mock.Mock
}

// Discover is a mocked method.
func (m *Mock) Discover() (lst kubernetes.MinionList, err error) {
	args := m.Called()
	x := args.Get(0)
	if x != nil {
		lst = x.(kubernetes.MinionList)
	}
	err = args.Error(1)
	return
}

// GetPodInfo is a mocked method.
func (m *Mock) GetPodInfo(namespace string, name string) (podInfo *kubernetes.PodInfo, err error) {
	args := m.Called(namespace, name)
	x := args.Get(0)
	if x != nil {
		podInfo = x.(*kubernetes.PodInfo)
	}
	err = args.Error(1)
	return
}
