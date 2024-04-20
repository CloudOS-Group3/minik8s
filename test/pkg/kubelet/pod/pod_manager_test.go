package pod

import (
	"minik8s/pkg/api"
	"minik8s/pkg/kubelet/pod"
	"testing"
)

func TestCreatePod(t *testing.T) {
	// create pod manager
	pm := pod.NewPodManager()

	// create pod
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

	if pm.CreatePod(newPod) == false {
		t.Fatalf("Failed to create pod")
	}
	pm.ShowPodInfo("test-pod")
}
