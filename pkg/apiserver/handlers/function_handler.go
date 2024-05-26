package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"net/http"
)

func AddFunction(context *gin.Context) {
	// Add function
	log.Info("Add function")
	var function api.Function
	if err := context.ShouldBind(&function); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}

	// Add function to etcd
	log.Info("Add function %s", function.Metadata.Name)

	URL := config.FunctionPath + function.Metadata.NameSpace + "/" + function.Metadata.Name
	functionByteArr, err := json.Marshal(function)
	if err != nil {
		log.Error("Error marshal function: %s", err.Error())
		return
	}
	etcdClient.PutEtcdPair(URL, string(functionByteArr))
}

func GetFunction(context *gin.Context) {
	// Get function
	log.Info("Get function")
	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)
	URL := config.FunctionPath + namespace + "/" + name
	function := etcdClient.GetEtcdPair(URL)
	var function_ api.Function
	if len(function) == 0 {
		log.Info("Function %s not found", name)
	} else {
		err := json.Unmarshal([]byte(function), &function_)
		if err != nil {
			log.Error("Error unmarshalling function json %v", err)
			return
		}
	}
	byteArr, err := json.Marshal(function_)
	if err != nil {
		log.Error("Error marshal function: %s", err.Error())
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"data": string(byteArr),
	})
}

func UpdateFunction(context *gin.Context) {
	// Update function
	log.Info("Update function")
	var function api.Function
	if err := context.ShouldBind(&function); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}

	URL := config.FunctionPath + function.Metadata.NameSpace + "/" + function.Metadata.Name
	functionByteArr, err := json.Marshal(function)
	if err != nil {
		log.Error("Error marshal function: %s", err.Error())
		return
	}
	etcdClient.PutEtcdPair(URL, string(functionByteArr))
}

func DeleteFunction(context *gin.Context) {
	// Delete function
	log.Info("Delete function")
	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)
	URL := config.FunctionPath + namespace + "/" + name
	etcdClient.DeleteEtcdPair(URL)
}
