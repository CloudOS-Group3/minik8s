package main

import "minik8s/pkg/kubeproxy"

func main() {
	proxy := kubeproxy.NewKubeProxy()
	proxy.Run()
}
