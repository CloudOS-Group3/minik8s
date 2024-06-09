package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/gpu"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"net/http"
	"time"
)

func AddGpuFunc(context *gin.Context) {
	log.Info("Add gpu spec")
	var gpujob api.GPUJob
	if err := context.ShouldBind(&gpujob); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}
	gpujob.Metadata.UUID = stringutil.GenerateRandomString(5)

	// Step1: Build image
	err := gpu.CreateGpuImage(&gpujob)
	if err != nil {
		log.Error("Error create image: %s", err.Error())
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}

	// Step2: Create pod to run the image
	pod := gpu.CreateGPUPod(&gpujob)
	etcdClient.PutPod(*pod)
	msg := msg_type.PodMsg{
		Opt:    msg_type.Add,
		NewPod: *pod,
	}
	msg_json, _ := json.Marshal(msg)
	publisher.Publish(msg_type.PodTopic, string(msg_json))

	gpujob.StartTime = time.Now().Format("2006-01-02 15:04:05")

	// Step3: Save gpujob to etcd, we save as <name>-<uuid>
	URL := config.GPUjobPath + gpujob.Metadata.Name + "-" + gpujob.Metadata.UUID
	gpuByteArr, err := json.Marshal(gpujob)
	if err != nil {
		log.Error("Error marshal gpu spec: %s", err.Error())
		return
	}
	etcdClient.PutEtcdPair(URL, string(gpuByteArr))
}

func GetAllGpuJobs(context *gin.Context) {
	URL := config.GPUjobPath
	gpujobs := etcdClient.PrefixGet(URL)

	jsonString := stringutil.EtcdResEntryToJSON(gpujobs)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})

}

func GetGpuJobsByName(context *gin.Context) {
	name := context.Param("name")
	URL := config.GPUjobPath + name
	gpujobs := etcdClient.PrefixGet(URL)

	jsonString := stringutil.EtcdResEntryToJSON(gpujobs)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})

}
