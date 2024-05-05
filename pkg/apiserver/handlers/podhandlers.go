package handlers

import (
	"encoding/json"
	"fmt"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetPods(context *gin.Context) {
	log.Info("received get pods request")

	URL := config.EtcdPodPath
	log.Debug("before prefix get")
	pods := etcdClient.PrefixGet(URL)

	log.Debug("get all pods are: %+v", pods)

	var podString []string
	for _, pod := range pods {
		podString = append(podString, pod.Value)
	}
	jsonValue := strings.Join(podString, ",")
	jsonValue = fmt.Sprint("[", jsonValue ,"]")

	context.JSON(http.StatusOK, gin.H{
		"data": jsonValue,
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
	
	newPod.Status.StartTime = time.Now()
	newPod.Metadata.UUID = uuid.NewString()

	etcdClient.PutPod(newPod)

	podByteArray, err := json.Marshal(newPod)

	log.Debug("pod byte array is: %+v", podByteArray)
	if err != nil {
		log.Error("Error: json marshal failed")
		return
	}

	// publisher.Publish("pod", string(podByteArray))

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
