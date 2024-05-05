package controllermanager

import "minik8s/pkg/controller/controllers"

type ControllerManager struct {
	DeploymentController *controllers.DeploymentController
}

func NewControllerManager() *ControllerManager {
	newDC := &controllers.DeploymentController{}

	return &ControllerManager{
		DeploymentController: newDC,
	}
}

func (CM *ControllerManager) Run() {
	go CM.DeploymentController.Run()
}