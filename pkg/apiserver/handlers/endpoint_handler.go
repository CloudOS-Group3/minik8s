package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"net/http"
)

func getEndpoints(context *gin.Context) {
	namespace := context.Param(config.NamespaceParam)
	label := context.Param(config.LabelParam)
	URL := config.EndpointPath + namespace + "/" + label
	endpoints := etcdClient.PrefixGet(URL)

	context.JSON(http.StatusOK, gin.H{
		"data": endpoints,
	})
}
func addEndpoints(context *gin.Context) {
	label := context.Param(config.LabelParam)
	var newEndpoint api.EndPoint
	if err := context.ShouldBind(&newEndpoint); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}
	endpointByteArray, err := json.Marshal(newEndpoint)
	if err != nil {
		return
	}
	URL := config.EndpointPath + newEndpoint.NameSpace + "/" + label
	etcdClient.PutEtcdPair(URL, string(endpointByteArray))
	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func deleteEndpoints(context *gin.Context) {
	namespace := context.Param(config.NamespaceParam)
	label := context.Param(config.LabelParam)
	URL := config.EndpointPath + namespace + "/" + label
	etcdClient.DeleteEtcdPair(URL)
	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
