package util

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/pkg/config"
	"time"
)

func GetEtcdClient() (*clientv3.Client, error) {
	EndPointURL := "http://" + config.Remotehost + ":2379"
	return clientv3.New(
		clientv3.Config{
			Endpoints:   []string{EndPointURL},
			DialTimeout: 2 * time.Second,
		})
}
