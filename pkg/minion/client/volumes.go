/*
 * Copyright (C) 2016 wikiwi.io
 *
 * This software may be modified and distributed under the terms
 * of the MIT license. See the LICENSE file for details.
 */

package client

import (
	"fmt"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
)

// VolumesInterface is part of the Clients interface.
type VolumesInterface interface {
	Volumes() VolumesService
}

// VolumesService is an interface for interacting with volumes in the minion API.
type VolumesService interface {
	List(podUID string) (*api.VolumeList, error)
	Get(podUID string, name string) (*api.Volume, error)
	Freeze(podUID string, name string) (*api.Volume, error)
	Thaw(podUID string, name string) (*api.Volume, error)
}

var _ VolumesService = new(volumesServiceImpl)

type volumesServiceImpl struct {
	client *Client
}

func (v *volumesServiceImpl) List(podUID string) (*api.VolumeList, error) {
	req, err := v.client.NewRequest("GET", fmt.Sprintf("volumes/%s", podUID), nil)
	if err != nil {
		return nil, err
	}
	var volumeList api.VolumeList
	_, err = v.client.Do(req, &volumeList)
	return &volumeList, err
}

func (v *volumesServiceImpl) Get(podUID string, name string) (*api.Volume, error) {
	req, err := v.client.NewRequest("GET", fmt.Sprintf("volumes/%s/%s", podUID, name), nil)
	if err != nil {
		return nil, err
	}
	var volume api.Volume
	_, err = v.client.Do(req, &volume)
	return &volume, err
}

func (v *volumesServiceImpl) Freeze(podUID string, name string) (*api.Volume, error) {
	req, err := v.client.NewRequest("POST", fmt.Sprintf("volumes/%s/%s", podUID, name), api.FreezeThawRequest{Action: "freeze"})
	if err != nil {
		return nil, err
	}
	var volume api.Volume
	_, err = v.client.Do(req, &volume)
	return &volume, err
}

func (v *volumesServiceImpl) Thaw(podUID string, name string) (*api.Volume, error) {
	req, err := v.client.NewRequest("POST", fmt.Sprintf("volumes/%s/%s", podUID, name), api.FreezeThawRequest{Action: "thaw"})
	if err != nil {
		return nil, err
	}
	var volume api.Volume
	_, err = v.client.Do(req, &volume)
	return &volume, err
}
