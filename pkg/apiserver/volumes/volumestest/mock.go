/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package volumestest

import (
	"github.com/stretchr/testify/mock"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
)

// ManagerMock implements a Mock for Manager.
type ManagerMock struct {
	mock.Mock
}

// List is a mocked method.
func (m *ManagerMock) List(namespace, name string) (list *api.VolumeList, err error) {
	args := m.Called(namespace, name)
	x := args.Get(0)
	if x != nil {
		list = x.(*api.VolumeList)
	}
	err = args.Error(1)
	return
}

// Get is a mocked method.
func (m *ManagerMock) Get(namespace, name, volume string) (vol *api.Volume, err error) {
	args := m.Called(namespace, name, volume)
	x := args.Get(0)
	if x != nil {
		vol = x.(*api.Volume)
	}
	err = args.Error(1)
	return
}

// Freeze is a mocked method.
func (m *ManagerMock) Freeze(namespace, name, volume string) (vol *api.Volume, err error) {
	args := m.Called(namespace, name, volume)
	x := args.Get(0)
	if x != nil {
		vol = x.(*api.Volume)
	}
	err = args.Error(1)
	return
}

// Thaw is a mocked method.
func (m *ManagerMock) Thaw(namespace, name, volume string) (vol *api.Volume, err error) {
	args := m.Called(namespace, name, volume)
	x := args.Get(0)
	if x != nil {
		vol = x.(*api.Volume)
	}
	err = args.Error(1)
	return
}
