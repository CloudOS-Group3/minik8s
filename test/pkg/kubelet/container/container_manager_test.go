package container

import (
	"minik8s/pkg/api"
	"minik8s/pkg/kubelet/container"
	"testing"
)

func TestContainerManager(t *testing.T) {
	cm := container.NewContainerManager()
	container_ := cm.CreateContainer(api.Container{
		Name:            "test-container",
		Image:           "docker.io/library/nginx:latest",
		ImagePullPolicy: api.PullPolicyIfNotPresent,
	}, "test")
	if container_ == nil {
		t.Fatalf("Failed to create container")
		return
	}
	if cm.StartContainerById(container_.ID(), "test") == false {
		t.Fatalf("Failed to start container")
		return
	}
}
