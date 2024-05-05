package controllermanager

import "minik8s/pkg/controller/controllers"

type ControllerManager struct {
	DeploymentController *controllers.DeploymentController
	EndpointController   *controllers.EndPointController
}

func NewControllerManager() *ControllerManager {
	newDC := &controllers.DeploymentController{}
	newEC := controllers.NewEndPointController()

	return &ControllerManager{
		DeploymentController: newDC,
		EndpointController:   newEC,
	}
}

func (CM *ControllerManager) Run() {
	go CM.DeploymentController.Run()
	go CM.EndpointController.Run()
}
