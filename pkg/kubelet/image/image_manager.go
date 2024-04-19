package image_manager

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"log"
	"minik8s/pkg/api"
)

type ImageManager struct {
}

// PullImage pulls the image from the registry.
// reference: https://www.rectcircle.cn/posts/containerd-2-client-core-process/
func (im *ImageManager) PullImage(imageName string, pullPolicy string) containerd.Image {
	client, err := containerd.New("/run/containerd/containerd.sock") // default address
	if err != nil {
		log.Printf("Failed to create containerd client: %v", err)
		return nil
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), "default")
	switch pullPolicy {
	case api.PullPolicyAlways:
		// Always pull the image
		break
	case api.PullPolicyIfNotPresent:
		// Check if the image exists
		image, err := client.ImageService().Get(context.Background(), imageName)
		if err == nil {
			log.Printf("Image %s already exists", imageName)
			return containerd.NewImage(client, image)
		}
	case api.PullPolicyNever:
		return nil
	default:
		break
	}
	image, err := client.Pull(ctx, imageName, containerd.WithPullUnpack)
	if err != nil {
		log.Printf("Failed to pull image %s: %v", imageName, err.Error())
		return nil
	}

	log.Printf("Image %s pulled successfully", image.Name())

	return image
}
