package etcd

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"minik8s/pkg/util"
	"minik8s/util/log"
	"sync"
)

type Store struct {
	etcdClient *clientv3.Client
}

func NewStore() *Store {
	cli, err := util.GetEtcdClient()
	if err != nil {
		log.Fatal("Failed to connect to etcd: %v", err.Error())
		return nil
	}

	return &Store{etcdClient: cli}
}

func (store Store) PutEtcdPair(key, value string) bool {
	// get an interface to access kv store
	kv := clientv3.NewKV(store.etcdClient)

	// then put the pair to kv store
	_, err := kv.Put(context.TODO(), key, value, clientv3.WithPrevKV())
	if err != nil {
		log.Fatal("Failed to put pair in kv-store: %v", err.Error())
		return false
	}

	return true
}

func (store Store) GetEtcdPair(key string) (value string) {
	kv := clientv3.NewKV(store.etcdClient)

	resp, err := kv.Get(context.TODO(), key)
	if err != nil {
		log.Fatal("Failed to get pair in kv-store: %v", err.Error())
		return ""
	}

	if len(resp.Kvs) == 0 {
		log.Warn("key %s not found in kv-store", key)
		return ""
	}

	return string(resp.Kvs[0].Value)
}

func (store Store) DeleteEtcdPair(key string) bool {
	kv := clientv3.NewKV(store.etcdClient)

	_, err := kv.Delete(context.TODO(), key)
	if err != nil {
		log.Fatal("Failed to delete pair in kv-store: %v", err.Error())
		return false
	}

	return true
}

func (store Store) PrefixWatch(wg *sync.WaitGroup, ctx context.Context, prefix string, handler func(key string, value string)) {
	/* usage (not pretty sure):
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	PrefixWatch(...)
	...
	cancel() (when you want to terminate it)
	wg.Wait()
	*/
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			default:
				rch := store.etcdClient.Watch(ctx, prefix, clientv3.WithPrefix())
				for resp := range rch {
					err := resp.Err()
					if err != nil {
						log.Fatal("Failed to watch prefix-watch: %s", err.Error())
					}
					for _, ev := range resp.Events {
						handler(string(ev.Kv.Key), string(ev.Kv.Value))
					}
				}
			}
		}
	}()
}
