package handlers

import (
	"encoding/json"
	"minik8s/pkg/api"
	msg "minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"net/http"
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

	jsonString := stringutil.EtcdResEntryToJSON(pods)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
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
		return
	}
	log.Debug("new pod is: %+v", newPod)

	newPod.Status.StartTime = time.Now()
	newPod.Metadata.UUID = uuid.NewString()
	newPod.Status.Phase = "Pending"

	// check if the pod already exists
	oldPod, exited := etcdClient.GetPod(newPod.Metadata.NameSpace, newPod.Metadata.Name)

	// need to interact with etcd
	etcdClient.PutPod(newPod)

	// construct message
	var message msg.PodMsg
	if exited {
		message = msg.PodMsg{
			Opt:    msg.Update,
			OldPod: oldPod,
			NewPod: newPod,
		}
	} else {
		message = msg.PodMsg{
			Opt:    msg.Add,
			NewPod: newPod,
		}

	}
	msg_json, _ := json.Marshal(message)

	publisher.Publish(msg.PodTopic, string(msg_json))
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
		log.Warn("pod %s doesn't exist", name)
	}

	log.Debug("get pod is: %+v", pod)

	byteArr, err := json.Marshal(pod)
	if err != nil {
		log.Error("Error marshal pod: %s", err.Error())
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"data": string(byteArr),
	})
}

func DeletePods(context *gin.Context) {
	log.Info("received delete all request")
	URL := config.EtcdPodPath
	ok := etcdClient.PrefixDelete(URL)
	if !ok {
		log.Error("delete pods failed")
		return
	}
}

func UpdatePod(context *gin.Context) {
	log.Info("received update pod request")

	var newPod api.Pod
	if err := context.ShouldBind(&newPod); err != nil {
		log.Error("error decode new pod")
		return
	}
	// check if the pod already exists
	oldPod, exited := etcdClient.GetPod(newPod.Metadata.NameSpace, newPod.Metadata.Name)

	// todo: should use update pod
	etcdClient.UpdatePod(newPod)

	// construct message
	var message msg.PodMsg
	if exited {
		message = msg.PodMsg{
			Opt:    msg.Update,
			OldPod: oldPod,
			NewPod: newPod,
		}
	} else {
		message = msg.PodMsg{
			Opt:    msg.Add,
			NewPod: newPod,
		}
	}
	msg_json, _ := json.Marshal(message)
	publisher.Publish(msg.PodTopic, string(msg_json))

}

func DeletePod(context *gin.Context) {
	log.Info("received delete pod request")

	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)

	// check if the pod already exists
	oldPod, exited := etcdClient.GetPod(namespace, name)
	if !exited {
		log.Error("pod not exist")
		return
	}
	etcdClient.DeletePod(namespace, name)

	// construct message
	message := msg.PodMsg{
		Opt:    msg.Delete,
		OldPod: oldPod,
	}
	msg_json, _ := json.Marshal(message)
	publisher.Publish(msg.PodTopic, string(msg_json))
}
