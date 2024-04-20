package container

import (
	"minik8s/pkg/api"
	"minik8s/pkg/kubelet/container"
	"testing"
)

func TestContainerManager(t *testing.T) {
	cm := container.ContainerManager{}
	container_ := cm.CreateContainer(api.Container{
		Name:            "test-container",
		Image:           "docker.io/library/nginx:latest",
		ImagePullPolicy: api.PullPolicyIfNotPresent,
	})
	if container_ == nil {
		t.Fatalf("Failed to create container")
		return
	}
}
