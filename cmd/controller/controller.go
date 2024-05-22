package main

import "minik8s/pkg/controller/controllermanager"

func main() {
	Controllers := controllermanager.NewControllerManager()
	Controllers.Run(make(chan bool))
}
