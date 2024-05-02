package ipvs

import (
	libipvs "github.com/moby/ipvs"
	"golang.org/x/sys/unix"
	"minik8s/pkg/api"
	"minik8s/util/log"
	"net"
)

type IPVS interface {
	AddService(service *api.Service) error
	UpdateService(service *api.Service) error
	DeleteService(service *api.Service) error
}

type IpvsHandler struct {
}

func (i *IpvsHandler) AddService(service *api.Service) error {
	handle, err := libipvs.New("")
	if err != nil {
		return err
	}
	services, err := handle.GetServices()
	if err != nil {
		return err
	}
	defer handle.Close()
	// add service
	for _, port := range service.Spec.Ports {
		svc := &libipvs.Service{
			// BUG : net.ParseIP("ClusterIP") is <nil>
			Address:       net.ParseIP(service.Spec.Type),
			Port:          uint16(port.Port),
			Protocol:      unix.IPPROTO_TCP,
			AddressFamily: unix.AF_INET, //nl.FAMILY_V4
			SchedName:     libipvs.RoundRobin,
		}

		// check if service already exists
		var is_existed bool = false
		for _, existed_svc := range services {
			if existed_svc.Address.String() == svc.Address.String() && existed_svc.Port == svc.Port {
				is_existed = true
				break
			}
		}
		if is_existed {
			continue
		}

		if err := handle.NewService(svc); err != nil {
			log.Fatal("Failed to add service: %v", err)
			return err
		}
	}

	return nil
}

func (i *IpvsHandler) UpdateService(service *api.Service) error {
	return nil
}

func (i *IpvsHandler) DeleteService(service *api.Service) error {
	return nil
}
