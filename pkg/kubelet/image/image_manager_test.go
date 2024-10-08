package image

import (
	"minik8s/pkg/api"
	"minik8s/pkg/util"
	"testing"
)

func TestImageManager(t *testing.T) {
	client, _ := util.CreateClient()
	if client == nil {
		t.Fatalf("Failed to create containerd client")
		return
	}
	image_ := PullImage("docker.io/library/nginx:latest", api.PullPolicyIfNotPresent, client, "test")
	if image_ == nil {
		t.Fatalf("Failed to pull image")
		return
	}
}
