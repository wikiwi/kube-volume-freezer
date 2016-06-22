package volumestest

import (
	"github.com/stretchr/testify/mock"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
)

type ManagerMock struct {
	mock.Mock
}

func (m *ManagerMock) List(namespace, name string) (list *api.VolumeList, err error) {
	args := m.Called(namespace, name)
	x := args.Get(0)
	if x != nil {
		list = x.(*api.VolumeList)
	}
	err = args.Error(1)
	return
}

func (m *ManagerMock) Get(namespace, name, volume string) (vol *api.Volume, err error) {
	args := m.Called(namespace, name, volume)
	x := args.Get(0)
	if x != nil {
		vol = x.(*api.Volume)
	}
	err = args.Error(1)
	return
}

func (m *ManagerMock) Freeze(namespace, name, volume string) (vol *api.Volume, err error) {
	args := m.Called(namespace, name, volume)
	x := args.Get(0)
	if x != nil {
		vol = x.(*api.Volume)
	}
	err = args.Error(1)
	return
}

func (m *ManagerMock) Thaw(namespace, name, volume string) (vol *api.Volume, err error) {
	args := m.Called(namespace, name, volume)
	x := args.Get(0)
	if x != nil {
		vol = x.(*api.Volume)
	}
	err = args.Error(1)
	return
}
