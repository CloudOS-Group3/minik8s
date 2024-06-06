package main

import (
	"minik8s/pkg/apiserver/server"
	"minik8s/pkg/config"
)

func main() {

	host, port := config.GetHostAndPort()
	server := server.NewAPIserver(host, port)

	server.Run()
}
