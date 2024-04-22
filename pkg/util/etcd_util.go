package util

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

func GetEtcdClient() (*clientv3.Client, error) {
	return clientv3.New(
		clientv3.Config{
			Endpoints:   []string{"http://localhost:2379"},
			DialTimeout: 2 * time.Second,
		})
}
