package test

import (
	"minik8s/pkg/api"
	image_manager "minik8s/pkg/kubelet/image"
	"testing"
)

func TestImageManager(t *testing.T) {
	im := image_manager.ImageManager{}
	image := im.PullImage("docker.io/library/nginx:latest", api.PullPolicyIfNotPresent)
	if image == nil {
		t.Fatalf("Failed to pull image")
		return
	}
}
