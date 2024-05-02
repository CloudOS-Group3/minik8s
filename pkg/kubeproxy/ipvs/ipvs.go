package ipvs

import (
	"github.com/moby/ipvs"
	"minik8s/pkg/api"
)

type IPVS interface {
	AddService(service *api.Service) error
	UpdateService(service *api.Service) error
	DeleteService(service *api.Service) error
}

type ipvs_handler struct {
}

func (i *ipvs_handler) AddService(service *api.Service) error {
	handle, err := ipvs.New("")
	if err != nil {
		return err
	}
	svcs, err := handle.GetServices()
	if err != nil {
		return err
	}
	defer handle.Close()
	// add service

	return nil
}
