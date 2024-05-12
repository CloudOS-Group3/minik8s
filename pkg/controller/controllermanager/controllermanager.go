package controllermanager

import "minik8s/pkg/controller/controllers"

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

func (CM *ControllerManager) Run() {
	go CM.DeploymentController.Run()
	go CM.EndpointController.Run()
	go CM.HPACcntroller.Run()
}
