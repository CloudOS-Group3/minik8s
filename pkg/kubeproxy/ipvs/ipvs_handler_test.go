package ipvs

import (
	"minik8s/pkg/api"
	"minik8s/util/log"
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

	// 检查 IPVS 是否包含添加的 TCP 服务
	output, err := exec.Command("ipvsadm", "-Ln").CombinedOutput()
	if err != nil {
		err := exec.Command("apt", "install", "ipvsadm").Run()
		if err != nil {
			log.Fatal("Failed to install ipvsadm: %v", err)
		}
	}

	if !strings.Contains(string(output), "TCP  0.0.0.0:81 rr") {
		t.Error("Expected TCP 0.0.0.0:81 rr in ipvsadm output, but not found")
	}
}
