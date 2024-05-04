package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"net/http"
)

func GetService(context *gin.Context) {
	namespace := context.Param(config.NamespaceParam)
	name := context.Param(config.NameParam)
	URL := config.ServicePath + namespace + "/" + name
	svc := etcdClient.GetEtcdPair(URL)

	context.JSON(http.StatusOK, gin.H{
		"data": svc,
	})
}

func GetAllServices(context *gin.Context) {
	URL := config.ServicePath
	services := etcdClient.PrefixGet(URL)

	context.JSON(http.StatusOK, gin.H{
		"data": services,
	})
}

func GetServicesByNamespace(context *gin.Context) {
	namespace := context.Param(config.NamespaceParam)
	URL := config.ServicePath + namespace
	services := etcdClient.PrefixGet(URL)

	context.JSON(http.StatusOK, gin.H{
		"data": services,
	})
}

func AddService(context *gin.Context) {
	var newService api.Service
	if err := context.ShouldBind(&newService); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}

	serviceByteArray, err := json.Marshal(newService)
	if err != nil {
		return
	}

	URL := config.ServicePath + newService.Metadata.NameSpace + "/" + newService.Metadata.Name
	etcdClient.PutEtcdPair(URL, string(serviceByteArray))

	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func DeleteService(context *gin.Context) {
	namespace := context.Param(config.NamespaceParam)
	name := context.Param(config.NameParam)
	URL := config.ServicePath + namespace + "/" + name
	etcdClient.DeleteEtcdPair(URL)

	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
