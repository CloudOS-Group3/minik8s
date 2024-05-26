package handlers

import (
	"context"
	"minik8s/pkg/config"
	"minik8s/pkg/etcd"
	"minik8s/pkg/kafka"
	"sync"
)

var publisher kafka.Publisher
var etcdClient etcd.Store

func WatchHandler(key string, value string) {
	// "/trigger/function_namespace/function_name"
	function := key[9:]
	// TODO: get function namespace and name, generate URL
	URL := ""
	str := etcdClient.GetEtcdPair(URL)
	if str == "" {
		return
	}
	// TODO: check function and send message
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
