/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

// Package clienttest provides utilities for testing the Minion Client.
package clienttest

import (
	"github.com/stretchr/testify/mock"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/client"
)

var _ client.Interface = new(Mock)
var _ client.VolumesService = new(VolumesServiceMock)

// Mock implements a Client Mock.
type Mock struct {
	VolumesMock *VolumesServiceMock
}

// Volumes returns the VolumesService Mock.
func (m *Mock) Volumes() client.VolumesService {
	return m.VolumesMock
}

// FactoryFake implements a Factory Fake.
type FactoryFake struct {
	// Clients is a static map returned during New().
	// Key of map is the address associated with the client.
	Clients map[string]client.Interface
}

// New returns a Client from the static map.
func (f *FactoryFake) New(address string) (client.Interface, error) {
	c, _ := f.Clients[address]
	return c, nil
}

// NewMock reutrns a new instance of the Client Mock.
func NewMock() *Mock {
	return &Mock{VolumesMock: &VolumesServiceMock{}}
}

// VolumesServiceMock implements a VolumesService Mock.
type VolumesServiceMock struct {
	mock.Mock
}

// List is a mocked method.
func (v *VolumesServiceMock) List(podUID string) (list *api.VolumeList, err error) {
	args := v.Called(podUID)
	x := args.Get(0)
	if x != nil {
		list = x.(*api.VolumeList)
	}
	err = args.Error(1)
	return
}

// Get is a mocked method.
func (v *VolumesServiceMock) Get(podUID string, name string) (vol *api.Volume, err error) {
	args := v.Called(podUID, name)
	x := args.Get(0)
	if x != nil {
		vol = x.(*api.Volume)
	}
	err = args.Error(1)
	return
}

// Freeze is a mocked method.
func (v *VolumesServiceMock) Freeze(podUID string, name string) (vol *api.Volume, err error) {
	args := v.Called(podUID, name)
	x := args.Get(0)
	if x != nil {
		vol = x.(*api.Volume)
	}
	err = args.Error(1)
	return
}

// Thaw is a mocked method.
func (v *VolumesServiceMock) Thaw(podUID string, name string) (vol *api.Volume, err error) {
	args := v.Called(podUID, name)
	x := args.Get(0)
	if x != nil {
		vol = x.(*api.Volume)
	}
	err = args.Error(1)
	return
}
