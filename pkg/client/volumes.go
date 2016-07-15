package client

import (
	"fmt"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
)

// VolumesInterface is part of the clients Interface.
type VolumesInterface interface {
	Volumes() VolumesService
}

// VolumesService is an interface for interacting with Volumes in the Master API.
type VolumesService interface {
	List(namespace, pod string) (*api.VolumeList, error)
	Get(namespace, pod, volume string) (*api.Volume, error)
	Freeze(namespace, pod, volume string) (*api.Volume, error)
	Thaw(namespace, pod, volume string) (*api.Volume, error)
}

var _ VolumesService = new(volumesServiceImpl)

type volumesServiceImpl struct {
	client *Client
}

func (v *volumesServiceImpl) List(namespace, pod string) (*api.VolumeList, error) {
	req, err := v.client.NewRequest("GET", fmt.Sprintf("volumes/%s/%s", namespace, pod), nil)
	if err != nil {
		return nil, err
	}
	var volumeList api.VolumeList
	_, err = v.client.Do(req, &volumeList)
	return &volumeList, err
}

func (v *volumesServiceImpl) Get(namespace, pod, volume string) (*api.Volume, error) {
	req, err := v.client.NewRequest("GET", fmt.Sprintf("volumes/%s/%s/%s", namespace, pod, volume), nil)
	if err != nil {
		return nil, err
	}
	var vol api.Volume
	_, err = v.client.Do(req, &vol)
	return &vol, err
}

func (v *volumesServiceImpl) Freeze(namespace, pod, volume string) (*api.Volume, error) {
	req, err := v.client.NewRequest("POST", fmt.Sprintf("volumes/%s/%s/%s", namespace, pod, volume), api.FreezeThawRequest{Action: "freeze"})
	if err != nil {
		return nil, err
	}
	var vol api.Volume
	_, err = v.client.Do(req, &vol)
	return &vol, err
}

func (v *volumesServiceImpl) Thaw(namespace, pod, volume string) (*api.Volume, error) {
	req, err := v.client.NewRequest("POST", fmt.Sprintf("volumes/%s/%s/%s", namespace, pod, volume), api.FreezeThawRequest{Action: "thaw"})
	if err != nil {
		return nil, err
	}
	var vol api.Volume
	_, err = v.client.Do(req, &vol)
	return &vol, err
}
