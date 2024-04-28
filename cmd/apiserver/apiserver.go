package main

import (
	"minik8s/pkg/apiserver/server"
	"minik8s/pkg/apiserver/config"
)


func main() {
	
	host, port := config.GetHostAndPort()
	server := server.NewApiserver(host, port)

	server.Run()
}
