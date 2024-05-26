package server

import (
	"fmt"
	"minik8s/pkg/apiserver/handlers"
	"minik8s/pkg/config"

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
	server.router.DELETE(config.PodsURL, handlers.DeletePods)
	server.router.GET(config.PodURL, handlers.GetPod)
	server.router.PUT(config.PodURL, handlers.UpdatePod)
	server.router.DELETE(config.PodURL, handlers.DeletePod)

	server.router.GET(config.DeploymentsURL, handlers.GetDeployments)
	server.router.POST(config.DeploymentsURL, handlers.AddDeployment)
	server.router.GET(config.DeploymentURL, handlers.GetDeployment)
	server.router.PUT(config.DeploymentURL, handlers.UpdateDeployment)
	server.router.DELETE(config.DeploymentURL, handlers.DeleteDeployment)

	server.router.GET(config.HPAsURL, handlers.GetHPAs)
	server.router.POST(config.HPAsURL, handlers.AddHPA)
	server.router.GET(config.HPAURL, handlers.GetHPA)
	server.router.PUT(config.HPAURL, handlers.UpdateHPA)
	server.router.DELETE(config.HPAURL, handlers.DeleteHPA)

	server.router.GET(config.LabelIndexURL, handlers.GetLabelIndex)
	server.router.POST(config.LabelIndexURL, handlers.AddLabelIndex)
	server.router.DELETE(config.LabelIndexURL, handlers.DeleteLabelIndex)

	server.router.GET(config.ServicesURL, handlers.GetServicesByNamespace)
	server.router.GET(config.ServicesAllURL, handlers.GetAllServices)
	server.router.GET(config.ServiceURL, handlers.GetService)
	server.router.PUT(config.ServiceURL, handlers.AddService)
	server.router.POST(config.ServiceURL, handlers.AddService)
	server.router.DELETE(config.ServiceURL, handlers.DeleteService)

	server.router.GET(config.DNSsURL, handlers.GetDNSs)
	server.router.POST(config.DNSsURL, handlers.AddDNS)
	server.router.DELETE(config.DNSURL, handlers.DeleteDNS)
	server.router.PUT(config.DNSURL, handlers.UpdateDNS)
	server.router.GET(config.DNSURL, handlers.GetDNS)

	server.router.GET(config.FunctionURL, handlers.GetFunction)
	server.router.PUT(config.FunctionURL, handlers.UpdateFunction)
	server.router.POST(config.FunctionURL, handlers.AddFunction)
	server.router.DELETE(config.FunctionURL, handlers.DeleteFunction)
}
