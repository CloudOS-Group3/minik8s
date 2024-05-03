package handlers

import (
	"encoding/json"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPods(context *gin.Context) {
	log.Info("received get pods request")

	URL := config.EtcdPodPath
	pods := etcdClient.PrefixGet(URL)

	log.Debug("get all pods are: %+v", pods)
	context.JSON(http.StatusOK, gin.H{
		"data": pods,
	})
}

func AddPod(context *gin.Context) {
	log.Info("received add pod request")

	var newPod api.Pod
	if err := context.ShouldBind(&newPod); err != nil {
		log.Error("decode pod failed")
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
	}
	log.Debug("new pod is: %+v", newPod)

	// need to interact with etcd
	etcdClient.PutPod(newPod)

	podByteArray, err := json.Marshal(newPod)

	if err != nil {
		log.Error("Error: json marshal failed")
		return
	}

	publisher.Publish("pod", string(podByteArray))

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

	pod, ok := etcdClient.GetPod(namespace, name)

	if !ok {
		log.Error("get pod not ok")
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data": pod,
	})
}

func UpdatePod(context *gin.Context) {
	log.Info("received update pod request")

	var newPod api.Pod
	if err := context.ShouldBind(&newPod); err != nil {
		log.Error("error decode new pod")
		return
	}

	// todo: should use update pod
	etcdClient.UpdatePod(newPod)

}

func DeletePod(context *gin.Context) {
	log.Info("received delete pod request")

	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)
	etcdClient.DeletePod(namespace, name)
}
