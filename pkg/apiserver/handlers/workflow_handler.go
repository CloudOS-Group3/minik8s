package handlers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"net/http"
)

func GetWorkflow(context *gin.Context) {
	// Get workflow
	log.Info("Get workflow")
	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)
	URL := config.WorkflowPath + namespace + "/" + name
	workflow := etcdClient.GetEtcdPair(URL)
	var workflow_ api.Workflow
	if len(workflow) == 0 {
		log.Info("Workflow %s not found", name)
	} else {
		err := json.Unmarshal([]byte(workflow), &workflow_)
		if err != nil {
			log.Error("Error unmarshalling workflow json %v", err)
			return
		}
	}
	byteArr, err := json.Marshal(workflow_)
	if err != nil {
		log.Error("Error marshal workflow: %s", err.Error())
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"data": string(byteArr),
	})
}

func AddWorkflow(context *gin.Context) {
	// Add workflow
	log.Info("Add workflow")
	var workflow api.Workflow
	if err := context.ShouldBind(&workflow); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}

	// Add workflow to etcd
	log.Info("Add workflow %s", workflow.Metadata.Name)

	URL := config.WorkflowPath + workflow.Metadata.NameSpace + "/" + workflow.Metadata.Name
	workflowByteArr, err := json.Marshal(workflow)
	if err != nil {
		log.Error("Error marshal workflow: %s", err.Error())
		return
	}
	etcdClient.PutEtcdPair(URL, string(workflowByteArr))
}

func DeleteWorkflow(context *gin.Context) {
	// Delete workflow
	log.Info("Delete workflow")
	name := context.Param(config.NameParam)
	namespace := context.Param(config.NamespaceParam)
	URL := config.WorkflowPath + namespace + "/" + name
	etcdClient.DeleteEtcdPair(URL)
}

func UpdateWorkflow(context *gin.Context) {
	// Update workflow
	log.Info("Update workflow")
	var workflow api.Workflow
	if err := context.ShouldBind(&workflow); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"status": "wrong",
		})
		return
	}

	// Update workflow to etcd
	log.Info("Update workflow %s", workflow.Metadata.Name)

	URL := config.WorkflowPath + workflow.Metadata.NameSpace + "/" + workflow.Metadata.Name
	workflowByteArr, err := json.Marshal(workflow)
	if err != nil {
		log.Error("Error marshal workflow: %s", err.Error())
		return
	}
	etcdClient.PutEtcdPair(URL, string(workflowByteArr))
}
