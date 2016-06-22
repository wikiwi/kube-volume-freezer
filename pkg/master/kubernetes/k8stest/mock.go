package k8stest

import (
	"github.com/stretchr/testify/mock"

	"github.com/wikiwi/kube-volume-freezer/pkg/master/kubernetes"
)

var _ kubernetes.Service = new(Mock)

type Mock struct {
	mock.Mock
}

func (m *Mock) Discover() (lst kubernetes.MinionList, err error) {
	args := m.Called()
	x := args.Get(0)
	if x != nil {
		lst = x.(kubernetes.MinionList)
	}
	err = args.Error(1)
	return
}

func (m *Mock) GetPodInfo(namespace string, name string) (podInfo *kubernetes.PodInfo, err error) {
	args := m.Called(namespace, name)
	x := args.Get(0)
	if x != nil {
		podInfo = x.(*kubernetes.PodInfo)
	}
	err = args.Error(1)
	return
}
