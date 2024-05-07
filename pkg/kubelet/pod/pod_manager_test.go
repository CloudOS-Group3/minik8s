package pod

import (
	"minik8s/pkg/api"
	"testing"
)

func TestCreatePod(t *testing.T) {

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

	if CreatePod(newPod) == false {
		t.Fatalf("Failed to create pod")
	}
	metric, err := GetPodMetrics(newPod)
	if err != nil {
		t.Fatalf("Failed to get pod metrics")
	}
	t.Logf("Pod metrics: %v", metric)

	// remove pod
	if DeletePod(newPod) == false {
		t.Fatalf("Failed to remove pod")
	}
}
