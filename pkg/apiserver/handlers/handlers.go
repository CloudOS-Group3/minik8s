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
	// get info of all nodes
}

func AddNode(context *gin.Context) {
	// add a new node into etcd
}

func GetNode(context *gin.Context) {
	// get info of a node
}

func DeleteNode(context *gin.Context) {
	// delete a node
}

func PutNode(context *gin.Context) {
	// change the data of a node
}

func GetPods(context *gin.Context) {
	log.Info("received get pods request")
}

func GetPod(context *gin.Context) {
	log.Info("received get pod request")
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

func UpdatePod(context *gin.Context) {
	log.Info("received update pod request")
}

func DeletePod(context *gin.Context) {
	log.Info("received delete pod request")
}
