package etcd

import (
	"context"
	"fmt"
	"minik8s/pkg/api"
	"sync"
	"testing"
)

type TestEtcdHandler struct {
	firstEnter bool
	done       chan bool
}

func (h TestEtcdHandler) WatchHandler(key string, value string) {
	fmt.Println("key: ", key, "value: ", value)
	select {
	case <-h.done:
		fmt.Println("it's done but enter again")
	default:
		close(h.done)
	}
}

func TestEtcd(t *testing.T) {
	newPod := &api.Pod{
		Metadata: api.ObjectMeta{
			Name:      "test-pod",
			NameSpace: "default",
		},
		Spec: api.PodSpec{
			Containers: []api.Container{
				{
					Name:            "test-container1",
					Image:           "docker.io/library/nginx:latest",
					ImagePullPolicy: api.PullPolicyIfNotPresent,
				},
				{
					Name:            "test-container2",
					Image:           "docker.io/library/nginx:latest",
					ImagePullPolicy: api.PullPolicyIfNotPresent,
				},
			},
		},
	}

	res := EtcdStore.PutPod(*newPod)
	if res != true {
		t.Errorf("etcd put pod fail")
	}

	pod, res := EtcdStore.GetPod("default", "test-pod")
	if res != true {
		t.Errorf("etcd get pod fail")
	}
	if pod.Metadata.Name != "test-pod" && pod.Metadata.NameSpace != "default" {
		t.Errorf("etcd get pod fail")
	}

	res = EtcdStore.DeletePod("default", "test-pod")
	if res != true {
		t.Errorf("etcd delete pod fail")
	}
}

func TestEtcdWatch(t *testing.T) {
	newPod := &api.Pod{
		Metadata: api.ObjectMeta{
			Name:      "test-pod",
			NameSpace: "default",
		},
		Spec: api.PodSpec{
			Containers: []api.Container{
				{
					Name:            "test-container1",
					Image:           "docker.io/library/nginx:latest",
					ImagePullPolicy: api.PullPolicyIfNotPresent,
				},
				{
					Name:            "test-container2",
					Image:           "docker.io/library/nginx:latest",
					ImagePullPolicy: api.PullPolicyIfNotPresent,
				},
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	handler := TestEtcdHandler{done: make(chan bool), firstEnter: true}
	EtcdStore.PrefixWatch(wg, ctx, "/registry/pods/", handler.WatchHandler)
	time := 0
	for {
		select {
		case <-handler.done:
			fmt.Println("it's done")
			cancel()
			wg.Wait()
			_, res := EtcdStore.GetPod("default", "test-pod")
			if res == true {
				res = EtcdStore.DeletePod("default", "test-pod")
				if res != true {
					t.Errorf("etcd delete pod fail")
				}
			}
			return
		default:
			if time%2 == 0 {
				fmt.Println("start put")
				res := EtcdStore.PutPod(*newPod)
				if res != true {
					t.Errorf("etcd put pod fail")
				}
			} else {
				fmt.Println("start delete")
				res := EtcdStore.DeletePod("default", "test-pod")
				if res != true {
					t.Errorf("etcd delete pod fail")
				}
			}
			time = time + 1
		}
	}
}
