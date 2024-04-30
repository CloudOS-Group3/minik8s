package main

import "minik8s/pkg/scheduler"

func main() {
	server := scheduler.NewScheduler()
	server.Run()
}
