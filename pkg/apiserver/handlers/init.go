package handlers

import (
	"minik8s/pkg/config"
	"minik8s/pkg/etcd"
	"minik8s/pkg/kafka"
)

var publisher kafka.Publisher
var etcdClient etcd.Store

func init() {
	KafkaURL := config.Remotehost + ":9092"
	publisher = *kafka.NewPublisher([]string{KafkaURL})
	etcdClient = *etcd.NewStore()
}
