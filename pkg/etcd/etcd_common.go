package etcd

import (
	"context"
	"minik8s/pkg/util"
	"minik8s/util/log"
	"sync"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Store struct {
	etcdClient *clientv3.Client
}

type ResEntry struct {
	Key   string
	Value string
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

func (store Store) PrefixGet(prefix string) []ResEntry {
	log.Debug("before etcd get")
	response, err := store.etcdClient.Get(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		log.Error("error get prefix %s", prefix)
		return nil
	}
	ret := []ResEntry{}
	for _, kv := range response.Kvs {
		ret = append(ret, ResEntry{
			Key:   string(kv.Key),
			Value: string(kv.Value),
		})
	}
	log.Debug("return value is: %+v", ret)
	return ret
}

func (store Store) PrefixDelete(prefix string) bool {
	log.Info("before etcd delete")
	_, err := store.etcdClient.Delete(context.TODO(), prefix, clientv3.WithPrefix())
	if err != nil {
		log.Error("error delete prefix %s", prefix)
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
				log.Info("start watching prefix %s", prefix)
				rch := store.etcdClient.Watch(ctx, prefix, clientv3.WithPrefix())
				log.Info("watch done")
				for resp := range rch {
					log.Info("receive message %s", resp.Events[0].Kv.Key)
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
