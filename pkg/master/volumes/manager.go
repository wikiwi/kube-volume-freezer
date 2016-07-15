// Package volumes contains the business logic of the Volume Resource.
package volumes

import (
	k8serrors "k8s.io/kubernetes/pkg/api/errors"

	"github.com/wikiwi/kube-volume-freezer/pkg/api"
	"github.com/wikiwi/kube-volume-freezer/pkg/api/errors"
	"github.com/wikiwi/kube-volume-freezer/pkg/api/issues"
	"github.com/wikiwi/kube-volume-freezer/pkg/log"
	"github.com/wikiwi/kube-volume-freezer/pkg/master/kubernetes"
	"github.com/wikiwi/kube-volume-freezer/pkg/minion/client"
)

// Manager contains the buisness logic for the Volume Resource.
type Manager interface {
	List(namespace, pod string) (*api.VolumeList, error)
	Get(namespace, pod, volume string) (*api.Volume, error)
	Freeze(namespace, pod, volume string) (*api.Volume, error)
	Thaw(namespace, pod, volume string) (*api.Volume, error)
}

type manager struct {
	kubernetes    kubernetes.Service
	clientFactory client.Factory
}

// NewManager creates a new Volumes Manager.
func NewManager(k kubernetes.Service, cf client.Factory) Manager {
	return &manager{kubernetes: k, clientFactory: cf}
}

func (m *manager) List(namespace, pod string) (*api.VolumeList, error) {
	nfo, err := m.getPodInfo(namespace, pod)
	if err != nil {
		return nil, err
	}

	minionClient, err := m.retrieveClient(nfo.NodeName)
	if err != nil {
		return nil, err
	}

	return minionClient.Volumes().List(nfo.UID)
}

func (m *manager) Get(namespace, pod, volume string) (*api.Volume, error) {
	nfo, err := m.getPodInfo(namespace, pod)
	if err != nil {
		return nil, err
	}

	minionClient, err := m.retrieveClient(nfo.NodeName)
	if err != nil {
		return nil, err
	}

	return minionClient.Volumes().Get(nfo.UID, volume)
}

func (m *manager) Freeze(namespace, pod, volume string) (*api.Volume, error) {
	nfo, err := m.getPodInfo(namespace, pod)
	if err != nil {
		return nil, err
	}

	minionClient, err := m.retrieveClient(nfo.NodeName)
	if err != nil {
		return nil, err
	}

	return minionClient.Volumes().Freeze(nfo.UID, volume)
}

func (m *manager) Thaw(namespace, pod, volume string) (*api.Volume, error) {
	nfo, err := m.getPodInfo(namespace, pod)
	if err != nil {
		return nil, err
	}

	minionClient, err := m.retrieveClient(nfo.NodeName)
	if err != nil {
		return nil, err
	}

	return minionClient.Volumes().Thaw(nfo.UID, volume)
}

func (m *manager) retrieveClient(nodeName string) (client.Interface, error) {
	minion, err := m.getMinion(nodeName)
	if err != nil {
		return nil, err
	}

	return m.clientFactory(minion.Address)
}

func (m *manager) getPodInfo(namespace, pod string) (*kubernetes.PodInfo, error) {
	nfo, err := m.kubernetes.GetPodInfo(namespace, pod)
	if err != nil {
		if statusError, ok := err.(*k8serrors.StatusError); ok && k8serrors.IsNotFound(statusError) {
			return nil, errors.NotFound("Pod not found").
				Append(issues.PodNotFound("Pod \"%s/%s\" does not exist", namespace, pod))
		}
		return nil, err
	}
	return nfo, nil
}

func (m *manager) getMinion(nodeName string) (*kubernetes.Minion, error) {
	minions, err := m.kubernetes.Discover()
	if err != nil {
		return nil, err
	}

	var minion *kubernetes.Minion
	for _, search := range minions {
		if search.NodeName == nodeName {
			minion = search
			break
		}
	}

	if minion == nil {
		log.Instance().Errorf("Minion on Node %q not found", nodeName)
		return nil, errors.NotFound("Minion not found").
			Append(issues.MinionNotFound("Minion on Host %q was not found", nodeName))
	}

	return minion, nil
}
