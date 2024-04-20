package handlers

import (
	"minik8s/pkg/apiserver/config"
	"minik8s/util/log"

	"github.com/gin-gonic/gin"
)

// all of the following handlers need to call etcd

func GetNodes(context *gin.Context) {
	// get info of all nodes
}

func PostNodes(context *gin.Context) {
	// add a new node into etcd
}

func GetNode(context *gin.Context) {
	// get info of a node
	name := context.Param(config.NameParam)
	log.Info("name is: %s", name)
}

func DeleteNode(context *gin.Context) {
	// delete a node
}

func PutNode(context *gin.Context) {
	// change the data of a node
}