package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"net/http"
)

func AddTrigger(context *gin.Context) {
	// Add function
	log.Info("Add trigger")
	var trigger api.Trigger
	if err := context.ShouldBind(&trigger); err != nil {
		log.Info("trigger unmarshal error")
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}

	funcURL := config.FunctionPath + trigger.Spec.FunctionNamespace + "/" + trigger.Spec.FunctionName
	str := etcdClient.GetEtcdPair(funcURL)
	if str == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "unknown function",
		})
		return
	}

	var function api.Function
	_ = json.Unmarshal([]byte(str), &function)
	function.Trigger.Event = true
	byteArr, _ := json.Marshal(function)
	etcdClient.PutEtcdPair(funcURL, string(byteArr))
	TriggerURL := config.EtcdTriggerPath + trigger.Spec.FunctionNamespace + "/" + trigger.Spec.FunctionNamespace
	byteArr, _ = json.Marshal(trigger)
	etcdClient.PutEtcdPair(TriggerURL, string(byteArr))

}

func GetTriggers(context *gin.Context) {
	// Get function
	log.Info("Get triggers")
	URL := config.EtcdTriggerPath
	triggers := etcdClient.PrefixGet(URL)

	log.Debug("get all nodes are: %+v", triggers)
	jsonString := stringutil.EtcdResEntryToJSON(triggers)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})
}

func DeleteTrigger(context *gin.Context) {
	// Delete function
	log.Info("Delete function")
	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)
	TriggerURL := config.EtcdTriggerPath + namespace + "/" + name
	str := etcdClient.GetEtcdPair(TriggerURL)
	if str == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}
	var trigger api.Trigger
	_ = json.Unmarshal([]byte(str), &trigger)
	funcURL := config.FunctionPath + trigger.Spec.FunctionNamespace + "/" + trigger.Spec.FunctionNamespace
	str = etcdClient.GetEtcdPair(funcURL)
	if str == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "unknown function",
		})
		return
	}

	var function api.Function
	_ = json.Unmarshal([]byte(str), &function)
	function.Trigger.Event = false
	byteArr, _ := json.Marshal(function)
	etcdClient.PutEtcdPair(funcURL, string(byteArr))
	etcdClient.DeleteEtcdPair(TriggerURL)
}

func HttpTriggerFunction(context *gin.Context) {
	log.Info("Http trigger function")
	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)
	funcURL := config.FunctionPath + namespace + "/" + name
	str := etcdClient.GetEtcdPair(funcURL)
	if str == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "unknown function",
		})
		return
	}
	var function api.Function
	_ = json.Unmarshal([]byte(str), &function)
	log.Info("function: %v", function)
	if function.Trigger.Http == true {
		var msg msg_type.TriggerMsg
		msg.Function = function
		msg.UUID = uuid.NewString()
		if err := context.ShouldBind(&msg); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"status": err.Error(),
			})
			return
		}
		jsonString, _ := json.Marshal(msg)
		publisher.Publish(msg_type.TriggerTopic, string(jsonString))
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "function doesn't allow http trigger",
		})
		return
	}
}
