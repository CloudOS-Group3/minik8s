package ipvs

import (
	"fmt"
	libipvs "github.com/moby/ipvs"
	"golang.org/x/sys/unix"
	"minik8s/pkg/api"
	"minik8s/util/log"
	"net"
	"os/exec"
	"strconv"
)

type IPVS interface {
	AddService(service *api.Service) error
	UpdateService(service *api.Service) error
	DeleteService(service *api.Service) error
}

type IpvsHandler struct {
}

func AddService(service *api.Service) error {
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
			Address:       net.ParseIP(service.Status.ClusterIP),
			Port:          uint16(port.Port),
			Protocol:      unix.IPPROTO_TCP,
			AddressFamily: unix.AF_INET, //nl.FAMILY_V4
			SchedName:     libipvs.RoundRobin,
		}

		// check if service already exists
		var is_existed bool = false
		for _, existed_svc := range services {
			if existed_svc.Port == svc.Port && existed_svc.Address.String() == svc.Address.String() {
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
		// ip addr add
		cmd := exec.Command("ip", "addr", "add", service.Status.ClusterIP, "dev", "flannel.1")
		output, err := cmd.CombinedOutput()
		log.Info("cmd: %v", cmd)
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Command output:", string(output))
		}
	}

	return nil
}

func UpdateService(service *api.Service) error {
	return AddService(service)
}

func AddEndpoint(service *api.Service) error {
	log.Info("AddEndpoint, service: %v", service)
	for _, endpoint := range service.Status.EndPoints {
		log.Info("endpoint: %v", endpoint)
		for _, port := range endpoint.Ports {
			cmd := exec.Command("ipvsadm", "-a", "-t", service.Status.ClusterIP+":"+endpoint.ServicePort, "-r", endpoint.IP+":"+strconv.Itoa(int(port.ContainerPort)), "-m")
			log.Info("cmd: %v", cmd)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Error:", err)
				fmt.Println("Command output:", string(output))
			} else {
				fmt.Println("Command succeeded with output:", string(output))
			}
			log.Info("bind endpoint %s:%d to service %s:%d", endpoint.IP, port.ContainerPort, service.Status.ClusterIP, endpoint.ServicePort)
		}
	}
	return nil
}

func DeleteService(service *api.Service) error {
	err := DeleteEndpoint(&service.Status.EndPoints, service.Status.ClusterIP)
	if err != nil {
		return err
	}
	// ipvsadm -D -t <ClusterIP>:<Port>
	for _, port := range service.Spec.Ports {
		cmd := exec.Command("ipvsadm", "-D", "-t", service.Status.ClusterIP+":"+strconv.Itoa(int(port.Port)))
		log.Info("cmd: %v", cmd)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Error:", err)
			fmt.Println("Command output:", string(output))
		}
	}
	// ip addr del
	cmd := exec.Command("ip", "addr", "del", service.Status.ClusterIP, "dev", "flannel.1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		fmt.Println("Command output:", string(output))
	}

	return nil
}

func UpdateEndpoints(a *api.Service, old *[]api.EndPoint) error {
	// delete old endpoints
	err := DeleteEndpoint(old, a.Status.ClusterIP)
	if err != nil {
		return err
	}
	// add new endpoints
	return AddEndpoint(a)
}

func DeleteEndpoint(old *[]api.EndPoint, ClusterIp string) error {
	for _, endpoint := range *old {
		for _, port := range endpoint.Ports {
			cmd := exec.Command("ipvsadm", "-d", "-t", ClusterIp+":"+endpoint.ServicePort, "-r", endpoint.IP+":"+strconv.Itoa(int(port.ContainerPort)))
			log.Info("cmd: %v", cmd)
			output, err := cmd.CombinedOutput()
			if err != nil {
				fmt.Println("Error:", err)
				fmt.Println("Command output:", string(output))
			} else {
				fmt.Println("Command succeeded with output:", string(output))
			}
			log.Info("unbind endpoint %s:%d from service %s:%d", endpoint.IP, port.ContainerPort, ClusterIp, endpoint.ServicePort)
		}
	}
	return nil
}
