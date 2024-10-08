package handlers

import (
	"encoding/json"
	"minik8s/pkg/api"
	msg "minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/util/consul"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetNodes(context *gin.Context) {
	log.Info("received get nodes request")

	URL := config.EtcdNodePath
	nodes := etcdClient.PrefixGet(URL)

	log.Debug("get all nodes are: %+v", nodes)
	jsonString := stringutil.EtcdResEntryToJSON(nodes)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
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
	newNode.Status.Condition.Status = api.NodeUnknown
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

	ID := "node-exporter-" + newNode.Metadata.Name
	name := "node-exporter-" + newNode.Metadata.Name
	port := 9100
	addr := newNode.Spec.NodeIP
	consul.RegisterService(ID, name, addr, port)

}

func GetNode(context *gin.Context) {
	log.Info("received get node request")
	name := context.Param(config.NameParam)

	if name == "" {
		log.Error("node name empty")
		return
	}

	URL := config.EtcdNodePath + name
	nodeJson := etcdClient.GetEtcdPair(URL)

	var node api.Node
	json.Unmarshal([]byte(nodeJson), &node)

	log.Info("node info: %+v", node)

	byteArr, err := json.Marshal(node)

	if err != nil {
		log.Error("error json marshal node: %s", err.Error())
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data": string(byteArr),
	})
}

func DeleteNode(context *gin.Context) {
	log.Info("received delete node request")
	name := context.Param(config.NameParam)

	URL := config.EtcdNodePath + name
	oldNode := etcdClient.GetEtcdPair(URL)
	etcdClient.DeleteEtcdPair(URL)
	ID := "node-exporter-" + name
	consul.DeRegisterService(ID)
	var node api.Node
	_ = json.Unmarshal([]byte(oldNode), &node)
	var message msg.NodeMsg
	message = msg.NodeMsg{
		Opt:     msg.Delete,
		OldNode: node,
	}
	msgJson, _ := json.Marshal(message)
	publisher.Publish(msg.NodeTopic, string(msgJson))
	for _, pod := range node.Status.Pods {
		pod.Spec.NodeName = ""
		URL = config.GetUrlPrefix() + config.PodURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, pod.Metadata.NameSpace, -1)
		URL = strings.Replace(URL, config.NamePlaceholder, pod.Metadata.Name, -1)
		podJson, _ := json.Marshal(pod)
		httputil.Put(URL, podJson)
	}
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
	log.Info("node info: %+v", newNode)
	nodeByteArray, err := json.Marshal(newNode)

	if err != nil {
		log.Error("error marshal newNode to json string")
	}

	URL := config.EtcdNodePath + newNode.Metadata.Name
	oldNode := etcdClient.GetEtcdPair(URL)
	etcdClient.PutEtcdPair(URL, string(nodeByteArray))

	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

	var message msg.NodeMsg
	if oldNode == "" {
		message = msg.NodeMsg{
			Opt:     msg.Add,
			NewNode: newNode,
		}
		ID := "node-exporter-" + newNode.Metadata.Name
		name := "node-exporter-" + newNode.Metadata.Name
		port := 9100
		addr := newNode.Spec.NodeIP
		consul.RegisterService(ID, name, addr, port)
	} else {
		var node api.Node
		if err := json.Unmarshal([]byte(oldNode), &node); err != nil {
			log.Error("error unmarshal old node")
		}
		message = msg.NodeMsg{
			Opt:     msg.Update,
			OldNode: node,
			NewNode: newNode,
		}
	}
	msg_json, _ := json.Marshal(message)
	publisher.Publish(msg.NodeTopic, string(msg_json))
}
