package handlers

import (
	"minik8s/pkg/api"
	"minik8s/pkg/apiserver/config"
	"minik8s/util/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// all of the following handlers need to call etcd

func GetNodes(context *gin.Context) {
	log.Info("received get nodes request")
}

func AddNode(context *gin.Context) {
	log.Info("received add node request")
}

func GetNode(context *gin.Context) {
	log.Info("received get node request")
	name := context.Param(config.NameParam)

	if name == "" {
		log.Error("node name empty")
		return
	}

	// todo: should query from etcd here
	// and we dont have node object yet...

	context.JSON(http.StatusOK, gin.H{
		"data": "",
	})
}

func DeleteNode(context *gin.Context) {
	log.Info("received delete node request")
}

func UpdateNode(context *gin.Context) {
	log.Info("received udpate node request")
}

func GetPods(context *gin.Context) {
	log.Info("received get pods request")
}

func AddPod(context *gin.Context) {
	// for the time being we dont interact with etcd
	// instead we directly communicate with pod manager
	log.Info("received add pod request")

	newPod := &api.Pod{}
	if err := context.ShouldBind(newPod); err != nil {
		log.Error("decode pod failed")
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
	}
	log.Debug("new pod is: %+v", newPod)

	// need to interact with etcd

}

func GetPod(context *gin.Context) {
	log.Info("received get pod request")

	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)

	if name == "" {
		log.Error("pod name empty")
		return
	}

	if namespace == "" {
		log.Error("namespace empty")
		return
	}

	pod := &api.Pod{}
	// todo: should query from etcd here

	context.JSON(http.StatusOK, gin.H{
		"data": *pod,
	})
}

func UpdatePod(context *gin.Context) {
	log.Info("received update pod request")
}

func DeletePod(context *gin.Context) {
	log.Info("received delete pod request")
}
