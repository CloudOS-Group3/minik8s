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

func NewAPIserver(host string, port int) *apiServer {
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
	server.router.PUT(config.ServiceURL, handlers.UpdateService)
	server.router.POST(config.ServiceURL, handlers.AddService)
	server.router.DELETE(config.ServiceURL, handlers.DeleteService)

	server.router.GET(config.DNSsURL, handlers.GetDNSs)
	server.router.POST(config.DNSsURL, handlers.AddDNS)
	server.router.DELETE(config.DNSURL, handlers.DeleteDNS)
	server.router.PUT(config.DNSURL, handlers.UpdateDNS)
	server.router.GET(config.DNSURL, handlers.GetDNS)

	server.router.POST(config.PersistentVolumesURL, handlers.AddPV)
	server.router.GET(config.PersistentVolumesURL, handlers.GetPVs)
	server.router.GET(config.PersistentVolumeURL, handlers.GetPV)
	server.router.DELETE(config.PersistentVolumeURL, handlers.DeletePV)

	server.router.POST(config.PersistentVolumeClaimsURL, handlers.AddPVC)
	server.router.GET(config.PersistentVolumeClaimsURL, handlers.GetPVCs)
	server.router.GET(config.PersistentVolumeClaimURL, handlers.GetPVC)
	server.router.DELETE(config.PersistentVolumeClaimURL, handlers.DeletePVC)

	server.router.GET(config.FunctionURL, handlers.GetFunction)
	server.router.GET(config.FunctionsURL, handlers.GetFunctions)
	server.router.PUT(config.FunctionURL, handlers.UpdateFunction)
	server.router.POST(config.FunctionURL, handlers.AddFunction)
	server.router.DELETE(config.FunctionURL, handlers.DeleteFunction)

	server.router.GET(config.WorkflowURL, handlers.GetWorkflow)
	server.router.PUT(config.WorkflowURL, handlers.UpdateWorkflow)
	server.router.POST(config.WorkflowURL, handlers.AddWorkflow)
	server.router.DELETE(config.WorkflowURL, handlers.DeleteWorkflow)
	server.router.GET(config.WorkflowsURL, handlers.GetAllWorkflow)

	server.router.GET(config.TriggersURL, handlers.GetTriggers)
	server.router.POST(config.TriggersURL, handlers.AddTrigger)
	server.router.DELETE(config.TriggerURL, handlers.DeleteTrigger)
	server.router.DELETE(config.TriggerWorkflowURL, handlers.DeleteWorkflowTrigger)
	server.router.POST(config.FunctionRunURL, handlers.HttpTriggerFunction)
	server.router.POST(config.WorkflowRunURL, handlers.HttpTriggerWorkflow)

	server.router.PUT(config.TriggerResultURL, handlers.UpdateTriggerResult)
	server.router.GET(config.TriggerResultURL, handlers.GetTriggerResult)
	server.router.GET(config.TriggerResultsURL, handlers.GetTriggerResults)

	server.router.GET(config.JobsURL, handlers.GetJobs)
	server.router.POST(config.JobsURL, handlers.AddJob)
	server.router.GET(config.JobURL, handlers.GetJob)
	server.router.DELETE(config.JobURL, handlers.DeleteJob)
	server.router.PUT(config.JobURL, handlers.UpdateJob)
	server.router.POST(config.JobResultURL, handlers.JobResultHandler)
	server.router.POST(config.GpuJobResultURL, handlers.GpuResultHandler)

	server.router.POST(config.GPUJobURL, handlers.AddGpuFunc)
	server.router.GET(config.GPUJobsURL, handlers.GetAllGpuJobs)
	server.router.GET(config.GPUJobURL, handlers.GetGpuJobsByName)

}
