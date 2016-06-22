package kubernetes

import (
	"fmt"

	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/labels"
)

// Minion holds information about the discovered Minion.
type Minion struct {
	Address  string
	NodeName string
}

// MinionList is a list of discovered Minions.
type MinionList []*Minion

type PodInfo struct {
	UID      string
	NodeName string
}

type Service interface {
	Discover() (MinionList, error)
	GetPodInfo(namespace string, name string) (*PodInfo, error)
}

type DiscoveryConfig struct {
	Namespace string
	Selector  string
	Scheme    string
	Port      int
}

func NewService(client unversioned.Interface, cfg *DiscoveryConfig) (Service, error) {
	sel, err := labels.Parse(cfg.Selector)
	if err != nil {
		return nil, err
	}
	s := &svc{
		namespace: cfg.Namespace, selector: sel,
		scheme: cfg.Scheme, port: cfg.Port,
		client: client,
	}
	return s, nil
}

var _ Service = new(svc)

type svc struct {
	namespace string
	selector  labels.Selector
	scheme    string
	port      int
	client    unversioned.Interface
}

func (s *svc) Discover() (MinionList, error) {
	list := MinionList{}

	pods, err := s.client.Pods(s.namespace).List(api.ListOptions{LabelSelector: s.selector})
	if err != nil {
		return nil, err
	}
	for _, pod := range pods.Items {
		minion := &Minion{
			NodeName: pod.Spec.NodeName,
			Address:  fmt.Sprintf("%s://%s:%d", s.scheme, pod.Status.PodIP, s.port),
		}
		list = append(list, minion)
	}

	return list, nil
}

func (s *svc) GetPodInfo(namespace string, name string) (*PodInfo, error) {
	pod, err := s.client.Pods(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return &PodInfo{UID: string(pod.GetUID()), NodeName: pod.Spec.NodeName}, nil
}
