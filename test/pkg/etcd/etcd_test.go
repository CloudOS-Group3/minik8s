package etcd

import (
	"minik8s/pkg/api"
	"minik8s/pkg/etcd"
	"testing"
)

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

	res := etcd.EtcdStore.PutPod(*newPod)
	if res != true {
		t.Errorf("etcd put pod fail")
	}

	pod, res := etcd.EtcdStore.GetPod("default", "test-pod")
	if res != true {
		t.Errorf("etcd get pod fail")
	}
	if pod.Metadata.Name != "test-pod" && pod.Metadata.NameSpace != "default" {
		t.Errorf("etcd get pod fail")
	}

	res = etcd.EtcdStore.DeletePod("default", "test-pod")
	if res != true {
		t.Errorf("etcd delete pod fail")
	}
}
