package ipvs

import "minik8s/pkg/api"

func main() {
	ipvs_handler := &ipvs_handler{}
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
	err := ipvs_handler.AddService(&service)
	if err != nil {
		return
	}
}

// README.md, example:
//import (
//	"github.com/moby/ipvs"
//	"log"
//)
//
//func main() {
//	handle, err := ipvs.New("")
//	if err != nil {
//		log.Fatalf("ipvs.New: %s", err)
//	}
//	svcs, err := handle.GetServices()
//	if err != nil {
//		log.Fatalf("handle.GetServices: %s", err)
//	}
//}
