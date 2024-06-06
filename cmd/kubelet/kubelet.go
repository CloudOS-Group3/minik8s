package main

import (
	"minik8s/pkg/kubelet/subscriber"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		server := subscriber.NewKubeletSubscriber("")
		server.Run()
	} else {
		if os.Args[1] == "--name" {
			server := subscriber.NewKubeletSubscriber(os.Args[2])
			server.Run()
		}
	}
}
