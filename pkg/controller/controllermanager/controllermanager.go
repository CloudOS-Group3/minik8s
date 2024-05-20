package controllermanager

import (
	"minik8s/pkg/controller/controllers"
	"minik8s/util/log"
)

type ControllerManager struct {
	DeploymentController *controllers.DeploymentController
	EndpointController   *controllers.EndPointController
	HPACcntroller        *controllers.HPAController
}

func NewControllerManager() *ControllerManager {
	newDC := &controllers.DeploymentController{}
	newEC := controllers.NewEndPointController()
	newHC := &controllers.HPAController{}

	return &ControllerManager{
		DeploymentController: newDC,
		EndpointController:   newEC,
		HPACcntroller:        newHC,
	}
}

func (CM *ControllerManager) Run(stop chan bool) {
	go CM.DeploymentController.Run()
	go CM.EndpointController.Run()
	go CM.HPACcntroller.Run()

	_, ok := <-stop
	if !ok {
		log.Debug("stop chan closed")
	}
	log.Debug("received stop signal")
}
