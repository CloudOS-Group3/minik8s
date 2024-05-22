package main

import (
	"minik8s/pkg/kubelet/subscriber"
)

func main() {
	server := subscriber.NewKubeletSubscriber()
	server.Run()
}
