package image

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"minik8s/pkg/api"
	"minik8s/util/log"
)

// PullImage pulls the image from the registry.
// reference: https://www.rectcircle.cn/posts/containerd-2-client-core-process/
func PullImage(imageName string, pullPolicy string, client *containerd.Client, namespace string) containerd.Image {
	if namespace == "" {
		namespace = "default"
	}
	ctx := namespaces.WithNamespace(context.Background(), namespace)
	switch pullPolicy {
	case api.PullPolicyAlways:
		// Always pull the image
		break
	case api.PullPolicyIfNotPresent:
		// Check if the image exists
		image, err := client.ImageService().Get(context.Background(), imageName)
		if err == nil {
			log.Error("Image %s already exists", imageName)
			return containerd.NewImage(client, image)
		}
	case api.PullPolicyNever:
		return nil
	default:
		break
	}
	image, err := client.Pull(ctx, imageName, containerd.WithPullUnpack)
	if err != nil {
		log.Error("Failed to pull image %s: %v", imageName, err.Error())
		return nil
	}

	log.Info("Image %s pulled successfully", image.Name())

	return image
}
