package handlers

import (
	"minik8s/pkg/etcd"
	"minik8s/pkg/kafka"
)


var publisher kafka.Publisher
var etcdClient etcd.Store

func init() {
	publisher = *kafka.NewPublisher([]string{"localhost:9092"})
	etcdClient = *etcd.NewStore()
}