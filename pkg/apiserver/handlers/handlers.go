package handlers

import (
	"encoding/json"
	"minik8s/pkg/api"
	"minik8s/pkg/apiserver/config"
	"minik8s/pkg/apiserver/serverconfig"
	"minik8s/pkg/etcd"
	"minik8s/pkg/kafka"
	"minik8s/util/log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// all of the following handlers need to call etcd

var publisher kafka.Publisher
var etcdClient etcd.Store

func init() {
	publisher = *kafka.NewPublisher([]string{"localhost:9092"})
	etcdClient = *etcd.NewStore()
}

func GetNodes(context *gin.Context) {
	log.Info("received get nodes request")

	URL := serverconfig.EtcdNodePath
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

	URL := serverconfig.EtcdNodePath + newNode.Metadata.Name
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

	URL := serverconfig.EtcdNodePath + name
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

	URL := serverconfig.EtcdNodePath + newNode.Metadata.Name
	etcdClient.PutEtcdPair(URL, string(nodeByteArray))

	context.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func GetPods(context *gin.Context) {
	log.Info("received get pods request")

	URL := serverconfig.EtcdPodPath
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
