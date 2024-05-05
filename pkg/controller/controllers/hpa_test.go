package controllers

import (
	"testing"
)

func TestGettingMachineInfo(t *testing.T) {
	hpa := &HPAController{}
	hpa.getContainerStatus()
}