package k8stest

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/api/errors"

	"github.com/wikiwi/kube-volume-freezer/pkg/master/kubernetes"
)

var _ kubernetes.Service = new(Fake)

type Fake struct {
	MinionList kubernetes.MinionList
	PodInfoMap map[string]*kubernetes.PodInfo // key is of format namespace/name
}

func (f *Fake) Discover() (kubernetes.MinionList, error) {
	return f.MinionList, nil
}

func (f *Fake) GetPodInfo(namespace string, name string) (*kubernetes.PodInfo, error) {
	key := namespace + "/" + name
	if nfo, found := f.PodInfoMap[key]; found {
		return nfo, nil
	}
	return nil, errors.NewNotFound(api.Resource("pods"), name)
}
