package main

import "minik8s/pkg/kubeproxy"

func main() {
	proxy := kubeproxy.NewKubeproxySub()
	proxy.Run()
}
