package server

import (
	"fmt"
	"minik8s/pkg/apiserver/config"
	"minik8s/pkg/apiserver/handlers"

	"github.com/gin-gonic/gin"
)

type apiServer struct {
	router *gin.Engine
	host   string
	port   int
}

func NewApiserver(host string, port int) *apiServer {
	return &apiServer{
		router: gin.Default(),
		host:   host,
		port:   port,
	}
}

func (server *apiServer) Run() {
	server.bind()

	URL := fmt.Sprintf("%s:%d", server.host, server.port)
	server.router.Run(URL)
}

// bind all urls to its handlers respectively
func (server *apiServer) bind() {

	server.router.GET(config.NodesURL, handlers.GetNodes)
	server.router.POST(config.NodesURL, handlers.AddNode)
	server.router.GET(config.NodeURL, handlers.GetNode)
	server.router.DELETE(config.NodeURL, handlers.DeleteNode)
	server.router.PUT(config.NodeURL, handlers.UpdateNode)

	server.router.GET(config.PodsURL, handlers.GetPods)
	server.router.POST(config.PodsURL, handlers.AddPod)
	server.router.GET(config.PodURL, handlers.GetPod)
	server.router.PUT(config.PodURL, handlers.UpdatePod)
	server.router.DELETE(config.PodURL, handlers.DeletePod)
}
