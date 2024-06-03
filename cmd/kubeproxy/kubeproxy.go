package main

import (
	"minik8s/pkg/kubeproxy"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		proxy := kubeproxy.NewKubeProxy("")
		proxy.Run()
	} else {
		if os.Args[1] == "--name" {
			proxy := kubeproxy.NewKubeProxy(os.Args[2])
			proxy.Run()
		}
	}

}
