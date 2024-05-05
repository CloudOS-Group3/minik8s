package controllers

import (
	"minik8s/util/log"

	"github.com/google/cadvisor/client"
)

type HPAController struct{}

func (hpa *HPAController) getContainerStatus() {
	client, err := client.NewClient("http://localhost:8080/")
	if err != nil {
		log.Error("Error create new cadvisor client")
		return
	}
	machineInfo, err := client.MachineInfo()
	if err != nil {
		log.Error("Error getting machine info: %s", err.Error())
		return
	}
	log.Info("machine info is: %+v", machineInfo)
}
