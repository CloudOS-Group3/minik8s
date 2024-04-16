package image_manager

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"log"
)

type ImageManager struct {
}

// PullImage pulls the image from the registry.
// reference: https://www.rectcircle.cn/posts/containerd-2-client-core-process/
func (im *ImageManager) pullImage(imageName string) error {
	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		return err
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "default")
	image, err := client.Pull(ctx, imageName, containerd.WithPullUnpack)
	if err != nil {
		return err
	}

	log.Printf("Image %s pulled successfully", image.Name())

	return nil
}
