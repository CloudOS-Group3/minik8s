package controllermanager

import (
	"minik8s/pkg/controller/controllers"
	"minik8s/util/log"
)

type ControllerManager struct {
	DeploymentController *controllers.DeploymentController
	EndpointController   *controllers.EndPointController
	HPAController        *controllers.HPAController
	NodeController       *controllers.NodeController
	ServerlessController *controllers.ServerlessController
	DNSController        *controllers.DNSController
	JobController        *controllers.JobController
}

func NewControllerManager() *ControllerManager {
	newDC := &controllers.DeploymentController{}
	newEC := controllers.NewEndPointController()
	newHC := &controllers.HPAController{}
	newNC := controllers.NewNodeController()
	newSC := controllers.NewServerlessController()
	newDNSController := controllers.NewDnsController()
	newJobController := controllers.NewJobController()

	return &ControllerManager{
		DeploymentController: newDC,
		EndpointController:   newEC,
		HPAController:        newHC,
		NodeController:       newNC,
		ServerlessController: newSC,
		DNSController:        newDNSController,
		JobController:        newJobController,
	}
}

func (CM *ControllerManager) Run(stop chan bool) {

	go CM.DeploymentController.Run()
	go CM.EndpointController.Run()
	go CM.HPAController.Run()
	go CM.NodeController.Run()
	go CM.DNSController.Run()
	go CM.JobController.Run()

	_, ok := <-stop
	if !ok {
		log.Debug("stop chan closed")
	}
	log.Debug("received stop signal")
}
