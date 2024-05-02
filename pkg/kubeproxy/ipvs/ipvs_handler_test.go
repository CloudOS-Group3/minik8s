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
					Port:       80,
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

	log.Info("ipvsadm output: %s", string(output))
	var ip = net.ParseIP(service.Spec.Type).String()
	if ip == "<nil>" {
		ip = "0.0.0.0"
	}
	var targetOutput = "TCP  " + ip + ":80 rr"
	if !strings.Contains(string(output), targetOutput) {
		t.Error("Expected "+targetOutput+" in ipvsadm output, but is %s", string(output))
	}
}
