package clienttest

import (
	"github.com/stretchr/testify/mock"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/client"
)

var _ client.Interface = new(Mock)
var _ client.VolumesService = new(VolumesServiceMock)

type Mock struct {
	VolumesMock *VolumesServiceMock
}

func (m *Mock) Volumes() client.VolumesService {
	return m.VolumesMock
}

type FactoryFake struct {
	Clients map[string]client.Interface
}

func (f *FactoryFake) New(address string) (client.Interface, error) {
	c, _ := f.Clients[address]
	return c, nil
}

func NewMock() *Mock {
	return &Mock{VolumesMock: &VolumesServiceMock{}}
}

type VolumesServiceMock struct {
	mock.Mock
}

func (v *VolumesServiceMock) List(podUID string) (list *api.VolumeList, err error) {
	args := v.Called(podUID)
	x := args.Get(0)
	if x != nil {
		list = x.(*api.VolumeList)
	}
	err = args.Error(1)
	return
}

func (v *VolumesServiceMock) Get(podUID string, name string) (vol *api.Volume, err error) {
	args := v.Called(podUID, name)
	x := args.Get(0)
	if x != nil {
		vol = x.(*api.Volume)
	}
	err = args.Error(1)
	return
}

func (v *VolumesServiceMock) Freeze(podUID string, name string) (vol *api.Volume, err error) {
	args := v.Called(podUID, name)
	x := args.Get(0)
	if x != nil {
		vol = x.(*api.Volume)
	}
	err = args.Error(1)
	return
}

func (v *VolumesServiceMock) Thaw(podUID string, name string) (vol *api.Volume, err error) {
	args := v.Called(podUID, name)
	x := args.Get(0)
	if x != nil {
		vol = x.(*api.Volume)
	}
	err = args.Error(1)
	return
}
