package kubeproxy

import "minik8s/pkg/api"

type ProxyInterface interface {
	OnServiceCreate(service *api.Service) error
	OnServiceUpdate(oldService, newService *api.Service) error
	OnServiceDelete(service *api.Service) error
}
