package image

import (
	"minik8s/pkg/api"
	"minik8s/pkg/kubelet/image"
	"minik8s/pkg/util"
	"testing"
)

func TestImageManager(t *testing.T) {
	im := image.ImageManager{}
	client, _ := util.CreateClient()
	if client == nil {
		t.Fatalf("Failed to create containerd client")
		return
	}
	image_ := im.PullImage("docker.io/library/nginx:latest", api.PullPolicyIfNotPresent, client, "test")
	if image_ == nil {
		t.Fatalf("Failed to pull image")
		return
	}
}
