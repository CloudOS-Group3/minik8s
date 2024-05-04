package service

import (
	"minik8s/pkg/api"
	"minik8s/pkg/etcd"
	"minik8s/util/log"
	"strconv"
)

const (
	etcdPrefix = "/registry/service"
)

type Service interface {
	AddService(service *api.Service) error
	UpdateService(service *api.Service) error
	DeleteService(service *api.Service) error
	GenerateClusterIP(service *api.Service) string
}
type ServiceManager struct {
	serviceList map[string]*api.Service // map namespace:name -> service
	Etcd        *etcd.Store
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		serviceList: make(map[string]*api.Service),
		Etcd:        etcd.NewStore(),
	}
}
func (s *ServiceManager) AddService(service *api.Service) error {
	var name = service.Metadata.NameSpace + ":" + service.Metadata.Name
	// check if existed
	if _, ok := s.serviceList[name]; ok {
		log.Info("Service %s already exists. Update it ...", name)
		return s.UpdateService(service)
	}

	// Generate ClusterIP
	service.Status.ClusterIP = s.GenerateClusterIP(service)

	// podList

	// Add service
	s.serviceList[name] = service
	return nil
}

func (s *ServiceManager) UpdateService(service *api.Service) error {
	return nil
}

func (s *ServiceManager) DeleteService(service *api.Service) error {
	return nil
}

func (s *ServiceManager) GenerateClusterIP(service *api.Service) string {
	// range: 172.16.0.0 - 172.16.0.255
	etcdURL := etcdPrefix + "/maxIP"
	maxIP := s.Etcd.GetEtcdPair(etcdURL)
	if len(maxIP) == 0 {
		s.Etcd.PutEtcdPair(etcdURL, "1")
		return "172.16.0.0"
	}
	i, _ := strconv.Atoi(maxIP)
	s.Etcd.PutEtcdPair(etcdURL, string(rune(i+1)))
	return "172.16.0." + maxIP
}

func (s *ServiceManager) GetService(namespace, name string) *api.Service {
	return s.serviceList[namespace+":"+name]
}
