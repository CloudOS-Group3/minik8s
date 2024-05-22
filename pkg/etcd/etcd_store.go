package etcd

var EtcdStore *Store = nil

func init() {
	// create an etcd store
	newStore := NewStore()
	EtcdStore = newStore
}
