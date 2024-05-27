package controllermanager

import (
	"minik8s/pkg/controller/controllers"
	"minik8s/util/log"
)

type ControllerManager struct {
	DeploymentController *controllers.DeploymentController
	EndpointController   *controllers.EndPointController
	HPACcntroller        *controllers.HPAController
	NodeController       *controllers.NodeController
	DNSController        *controllers.DNSController
	JobController        *controllers.JobController
}

func NewControllerManager() *ControllerManager {
	newDC := &controllers.DeploymentController{}
	newEC := controllers.NewEndPointController()
	newHC := &controllers.HPAController{}
	newNC := controllers.NewNodeController()
	newDNSController := controllers.NewDnsController()
	newJobController := controllers.NewJobController()

	return &ControllerManager{
		DeploymentController: newDC,
		EndpointController:   newEC,
		HPACcntroller:        newHC,
		NodeController:       newNC,
		DNSController:        newDNSController,
		JobController:        newJobController,
	}
}

func (CM *ControllerManager) Run(stop chan bool) {

	go CM.DeploymentController.Run()
	go CM.EndpointController.Run()
	go CM.HPACcntroller.Run()
	go CM.NodeController.Run()
	go CM.DNSController.Run()
	go CM.JobController.Run()

	_, ok := <-stop
	if !ok {
		log.Debug("stop chan closed")
	}
	log.Debug("received stop signal")
}
