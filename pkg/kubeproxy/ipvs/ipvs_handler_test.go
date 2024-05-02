package ipvs

import (
	"minik8s/pkg/api"
	"minik8s/util/log"
	"net"
	"os/exec"
	"strings"
	"testing"
)

func TestIpvsHandler_AddService(t *testing.T) {
	ipvsHandler := IpvsHandler{}
	service := api.Service{
		APIVersion: "v1",
		Kind:       "Service",
		Metadata: api.ObjectMeta{
			Name: "nginx",
		},
		Spec: api.ServiceSpec{
			Type: "ClusterIP",
			Ports: []api.ServicePort{
				{
					Port:       81,
					TargetPort: 8080,
					Protocol:   "TCP",
					Name:       "http",
				},
			},
		},
	}
	err := ipvsHandler.AddService(&service)
	if err != nil {
		t.Errorf("AddService() error = %v", err)
	}

	// check if service is added
	output, err := exec.Command("ipvsadm", "-Ln").CombinedOutput()
	if err != nil {
		err := exec.Command("apt", "install", "ipvsadm").Run()
		if err != nil {
			log.Fatal("Failed to install ipvsadm: %v", err)
		}
	}

	var targetOutput = "TCP  " + net.ParseIP(service.Spec.Type).String() + ":81 rr"
	if !strings.Contains(string(output), targetOutput) {
		t.Error("Expected "+targetOutput+" in ipvsadm output, but is %v", string(output))
	}
}
