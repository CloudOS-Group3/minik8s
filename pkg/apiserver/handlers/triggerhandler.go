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
	"time"
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

	if trigger.IsWorkflow {
		wfURL := config.WorkflowPath + trigger.Spec.FunctionNamespace + "/" + trigger.Spec.FunctionName
		str := etcdClient.GetEtcdPair(wfURL)
		if str == "" {
			context.JSON(http.StatusBadRequest, gin.H{
				"status": "unknown workflow",
			})
			return
		}

		var workflow api.Workflow
		_ = json.Unmarshal([]byte(str), &workflow)
		workflow.Trigger.Event = true
		byteArr, _ := json.Marshal(workflow)
		etcdClient.PutEtcdPair(wfURL, string(byteArr))
		TriggerURL := config.EtcdTriggerWorkflowPath + trigger.Spec.FunctionNamespace + "/" + trigger.Spec.FunctionNamespace
		byteArr, _ = json.Marshal(trigger)
		etcdClient.PutEtcdPair(TriggerURL, string(byteArr))
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
func DeleteWorkflowTrigger(context *gin.Context) {
	// Delete workflow trigger
	log.Info("Delete workflow trigger")
	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)
	TriggerURL := config.EtcdTriggerWorkflowPath + namespace + "/" + name
	str := etcdClient.GetEtcdPair(TriggerURL)
	if str == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}
	var trigger api.Trigger
	_ = json.Unmarshal([]byte(str), &trigger)
	wfURL := config.WorkflowPath + trigger.Spec.FunctionNamespace + "/" + trigger.Spec.FunctionNamespace
	str = etcdClient.GetEtcdPair(wfURL)
	if str == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "unknown workflow",
		})
		return
	}

	var workflow api.Workflow
	_ = json.Unmarshal([]byte(str), &workflow)
	workflow.Trigger.Event = false
	byteArr, _ := json.Marshal(workflow)
	etcdClient.PutEtcdPair(wfURL, string(byteArr))
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
func HttpTriggerWorkflow(context *gin.Context) {
	log.Info("Http trigger workflow")
	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)
	wfURL := config.WorkflowPath + namespace + "/" + name
	str := etcdClient.GetEtcdPair(wfURL)
	if str == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "unknown workflow",
		})
		return
	}
	var workflow api.Workflow
	_ = json.Unmarshal([]byte(str), &workflow)
	log.Info("workflow: %v", workflow)
	if workflow.Trigger.Http == true {
		var msg msg_type.WorkflowTriggerMsg
		msg.Workflow = workflow
		msg.UUID = uuid.NewString()
		if err := context.ShouldBind(&msg); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{
				"status": err.Error(),
			})
			return
		}
		log.Info("msg: %v", msg)
		jsonString, _ := json.Marshal(msg)
		publisher.Publish(msg_type.TriggerWorkflowTopic, string(jsonString))

		// store a empty result
		var result api.WorkflowResult
		result.Metadata = workflow.Metadata
		result.InvokeTime = time.Now().Format("2006-01-02 15:04:05")
		result.EndTime = "Running"
		URL := config.TriggerResultPath + msg.UUID
		byteArr, _ := json.Marshal(result)
		etcdClient.PutEtcdPair(URL, string(byteArr))
	} else {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrokflow doesn't allow http trigger",
		})
		return
	}
}

func UpdateTriggerResult(context *gin.Context) {
	// Add workflow
	log.Info("update workflow result")
	uuid := context.Param(config.UUIDParam)
	var res api.WorkflowResult
	if err := context.ShouldBind(&res); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}

	URL := config.TriggerResultPath + uuid
	workflowByteArr, err := json.Marshal(res)
	if err != nil {
		log.Error("Error marshal workflow: %s", err.Error())
		return
	}
	etcdClient.PutEtcdPair(URL, string(workflowByteArr))
}
func GetTriggerResult(context *gin.Context) {
	// Get function
	log.Info("Get trigger result")
	uuid := context.Param(config.UUIDParam)
	URL := config.TriggerResultPath + uuid
	str := etcdClient.GetEtcdPair(URL)
	if str == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"data": str,
	})
}
func GetTriggerResults(context *gin.Context) {
	// Get function
	log.Info("Get trigger results")
	URL := config.TriggerResultPath
	results := etcdClient.PrefixGet(URL)

	log.Debug("get all results are: %+v", results)
	jsonString := stringutil.EtcdResEntryToJSON(results)
	context.JSON(http.StatusOK, gin.H{
		"data": jsonString,
	})
}
