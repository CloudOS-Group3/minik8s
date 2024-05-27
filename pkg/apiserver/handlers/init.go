package handlers

import (
	"context"
	"encoding/json"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/etcd"
	"minik8s/pkg/kafka"
	"strings"
	"sync"
)

var publisher kafka.Publisher
var etcdClient etcd.Store

func WatchHandler(key string, value string) {
	// "/trigger/function_namespace/function_name"
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
		jsonString, _ := json.Marshal(msg)
		publisher.Publish(msg_type.TriggerTopic, string(jsonString))
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
		cancel()
		wg.Wait()
	}()
}
