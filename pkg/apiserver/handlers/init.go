package handlers

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/etcd"
	"minik8s/pkg/kafka"
	"minik8s/util/log"
	"strings"
	"sync"
)

var publisher kafka.Publisher
var etcdClient etcd.Store

func WatchHandler(key string, value string) {
	// "/trigger/function_namespace/function_name"
	log.Info("etcd trigger")
	str := key[9:]
	strList := strings.Split(str, "/")
	functionNamespace := strList[0]
	functionName := strList[1]
	URL := config.FunctionPath + functionNamespace + "/" + functionName
	str = etcdClient.GetEtcdPair(URL)
	if str == "" {
		return
	}
	var function api.Function
	_ = json.Unmarshal([]byte(str), &function)
	if function.Trigger.Event == true {
		var msg msg_type.TriggerMsg
		msg.Function = function
		msg.Params = value
		msg.UUID = uuid.NewString()
		jsonString, _ := json.Marshal(msg)
		publisher.Publish(msg_type.TriggerTopic, string(jsonString))
	}
}
func WatchWorkflowHandler(key string, value string) {
	// "/workflowtrigger/function_namespace/function_name"
	log.Info("etcd trigger")
	str := key[17:]
	strList := strings.Split(str, "/")
	wfNamespace := strList[0]
	wfName := strList[1]
	URL := config.WorkflowPath + wfNamespace + "/" + wfName
	str = etcdClient.GetEtcdPair(URL)
	if str == "" {
		return
	}
	var workflow api.Workflow
	_ = json.Unmarshal([]byte(str), &workflow)
	if workflow.Trigger.Event == true {
		var msg msg_type.WorkflowTriggerMsg
		msg.Workflow = workflow
		msg.Params = value
		msg.UUID = uuid.NewString()
		jsonString, _ := json.Marshal(msg)
		publisher.Publish(msg_type.TriggerWorkflowTopic, string(jsonString))
	}
}

func init() {
	KafkaURL := config.Remotehost + ":9092"
	publisher = *kafka.NewPublisher([]string{KafkaURL})
	etcdClient = *etcd.NewStore()
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		wg := &sync.WaitGroup{}
		prefix := config.UserTriggerPath
		etcdClient.PrefixWatch(wg, ctx, prefix, WatchHandler)
		prefix = config.UserTriggerWorkflowPath
		etcdClient.PrefixWatch(wg, ctx, prefix, WatchWorkflowHandler)
		wg.Wait()
		cancel()
	}()
}
