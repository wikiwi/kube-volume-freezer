package volumes

import (
	"path"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/api/errors"
	"github.com/wikiwi/kube-volume-freezer/pkg/api/issues"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/fs"
)

const (
	// PodsBasePath is the kubelets pod folder.
	PodsBasePath = "/var/lib/kubelet/pods"
)

type Manager interface {
	List(uid string) (*api.VolumeList, error)
	Get(uid, volume string) (*api.Volume, error)
	Freeze(uid, volume string) (*api.Volume, error)
	Thaw(uid, volume string) (*api.Volume, error)
}

type manager struct {
	fs fs.FileSystem
}

func NewManager(fs fs.FileSystem) Manager {
	return &manager{fs: fs}
}

func (m *manager) List(uid string) (*api.VolumeList, error) {
	ls, err := m.listVolumes(uid)
	if err != nil {
		return nil, err
	}
	return &api.VolumeList{PodUID: uid, Items: ls}, nil
}

func (m *manager) Get(uid, volume string) (*api.Volume, error) {
	_, err := m.resolveToAbsolutePath(uid, volume)
	if err != nil {
		return nil, err
	}
	return &api.Volume{PodUID: uid, Name: volume}, nil
}

func (m *manager) Freeze(uid, volume string) (*api.Volume, error) {
	volumePath, err := m.resolveToAbsolutePath(uid, volume)
	if err != nil {
		return nil, err
	}

	err = m.fs.Freeze(volumePath)
	if err != nil {
		return nil, err
	}

	return &api.Volume{PodUID: uid, Name: volume}, nil
}

func (m *manager) Thaw(uid, volume string) (*api.Volume, error) {
	volumePath, err := m.resolveToAbsolutePath(uid, volume)
	if err != nil {
		return nil, err
	}

	err = m.fs.Thaw(volumePath)
	if err != nil {
		return nil, err
	}

	return &api.Volume{PodUID: uid, Name: volume}, nil
}

// resolveToAbsoultePath returns the absolute path of the volume.
func (m *manager) resolveToAbsolutePath(uid, volume string) (string, error) {
	ls, err := m.listVolumesAbsolute(uid)
	if err != nil {
		return "", err
	}
	for _, x := range ls {
		if path.Base(x) == volume {
			return x, nil
		}
	}
	er := errors.NotFound("Volume not found").
		Append(issues.VolumeNotFound("Volume %q of Pod with UID %q does not exist", volume, uid))
	return "", er
}

// listVolumes returns all volumes of Pod. The returned values
// are the names of the volumes.
func (m *manager) listVolumes(uid string) ([]string, error) {
	ls, err := m.listVolumesAbsolute(uid)
	if err != nil {
		return nil, err
	}
	ret := []string{}
	for _, x := range ls {
		ret = append(ret, path.Base(x))
	}
	return ret, nil
}

// listVolumes returns all volumes of Pod. The returned values
// are absolute paths to the volume folder.
func (m *manager) listVolumesAbsolute(uid string) ([]string, error) {
	podVolumesFolder := PodsBasePath + "/" + uid + "/volumes"

	exists, err := m.fs.DirExists(podVolumesFolder)
	if err != nil {
		return nil, err
	}
	if !exists {
		er := errors.NotFound("Pod not found").
			Append(issues.PodNotFound("Pod with UID %q does not exist", uid))
		return nil, er
	}

	volumes := []string{}
	files, err := m.fs.ReadDir(podVolumesFolder)
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		if f.IsDir() {
			files2, err := m.fs.ReadDir(podVolumesFolder + "/" + f.Name())
			if err != nil {
				return nil, err
			}
			for _, f2 := range files2 {
				if f2.IsDir() {
					volumes = append(volumes, podVolumesFolder+"/"+f.Name()+"/"+f2.Name())
				}
			}
		}
	}
	return volumes, nil
}
