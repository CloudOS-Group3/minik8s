package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/api"
	msg "minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/controller/controllers"
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
	// check if the service already exists
	oldService, _ := controllers.GetService(newService.Metadata.NameSpace, newService.Metadata.Name)

	URL := config.ServicePath + newService.Metadata.NameSpace + "/" + newService.Metadata.Name
	etcdClient.PutEtcdPair(URL, string(serviceByteArray))

	//construct message
	var message msg.ServiceMsg
	if oldService != nil {
		message = msg.ServiceMsg{
			Opt:        msg.Update,
			OldService: *oldService,
			NewService: newService,
		}
	} else {
		message = msg.ServiceMsg{
			Opt:        msg.Add,
			NewService: newService,
		}
	}
	msg_json, _ := json.Marshal(message)
	publisher.Publish(msg.ServiceTopic, string(msg_json))

	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func DeleteService(context *gin.Context) {
	namespace := context.Param(config.NamespaceParam)
	name := context.Param(config.NameParam)
	URL := config.ServicePath + namespace + "/" + name

	// check if the service already exists
	oldService, _ := controllers.GetService(namespace, name)
	if oldService == nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}
	etcdClient.DeleteEtcdPair(URL)

	//construct message
	message := msg.ServiceMsg{
		Opt:        msg.Delete,
		OldService: *oldService,
	}
	msg_json, _ := json.Marshal(message)
	publisher.Publish(msg.ServiceTopic, string(msg_json))

	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
