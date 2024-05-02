package handlers

import (
	"encoding/json"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetNodes(context *gin.Context) {
	log.Info("received get nodes request")

	URL := config.EtcdNodePath
	nodes := etcdClient.PrefixGet(URL)

	log.Debug("get all nodes are: %+v", nodes)
	context.JSON(http.StatusOK, gin.H{
		"data": nodes,
	})
}

func AddNode(context *gin.Context) {
	log.Info("received add node request")

	var newNode api.Node
	if err := context.ShouldBind(&newNode); err != nil {
		log.Error("decode node failed")
		context.JSON(http.StatusOK, gin.H{
			"status": "wrong",
		})
		return
	}

	nodeByteArray, err := json.Marshal(newNode)

	if err != nil {
		log.Error("error marshal new node")
		return
	}

	URL := config.EtcdNodePath + newNode.Metadata.Name
	etcdClient.PutEtcdPair(URL, string(nodeByteArray))

	context.JSON(http.StatusOK, gin.H{
		"statas": "ok",
	})

}

func GetNode(context *gin.Context) {
	log.Info("received get node request")
	name := context.Param(config.NameParam)

	if name == "" {
		log.Error("node name empty")
		return
	}

	nodeJson := etcdClient.GetEtcdPair(name)

	var node api.Node
	json.Unmarshal([]byte(nodeJson), node)

	log.Info("node info: %+v", node)

	context.JSON(http.StatusOK, gin.H{
		"data": node,
	})
}

func DeleteNode(context *gin.Context) {
	log.Info("received delete node request")
	name := context.Param(config.NameParam)

	if name == "" {
		log.Error("node name empty")
		return
	}

	URL := config.EtcdNodePath + name
	etcdClient.DeleteEtcdPair(URL)
}

func UpdateNode(context *gin.Context) {
	log.Info("received udpate node request")

	var newNode api.Node
	if err := context.ShouldBind(&newNode); err != nil {
		log.Error("decode node failed")
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}

	nodeByteArray, err := json.Marshal(newNode)

	if err != nil {
		log.Error("error marshal newNode to json string")
	}

	URL := config.EtcdNodePath + newNode.Metadata.Name
	etcdClient.PutEtcdPair(URL, string(nodeByteArray))

	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
