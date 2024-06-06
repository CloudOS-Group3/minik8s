package main

import (
	"github.com/gin-gonic/gin"
	"minik8s/util/log"
)

var counter int

func count(context *gin.Context) {
	counter++
	context.JSON(200, gin.H{
		"counter": counter,
	})
}

func main() {
	counter = 0
	server := gin.Default()
	server.GET("/", count)
	err := server.Run(":12345")
	if err != nil {
		log.Fatal("server start error", err)
		return
	}
}
