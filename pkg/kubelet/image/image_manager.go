package image

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"minik8s/pkg/api"
	"minik8s/util/log"
	"os/exec"
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
	case api.PullFromRegistry:
		return pullFromRegistry(imageName, client, namespace)
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

func pullFromRegistry(imageName string, client *containerd.Client, namespace string) containerd.Image {
	cmd := exec.Command("nerdctl", "-n", namespace, "pull", imageName)
	log.Info("cmd: %v", cmd)
	output, err := cmd.CombinedOutput()
	log.Info("output: %s", string(output))
	if err != nil {
		log.Error("Failed to run nerdctl pull: %s", err.Error())
		return nil
	}
	ctx := namespaces.WithNamespace(context.Background(), namespace)
	//image, err := client.ImageService().Get(ctx, imageName)
	image, err := client.GetImage(ctx, imageName+":latest")
	if err != nil {
		log.Error("Failed to get image %s: %v", imageName, err.Error())
		return nil
	}
	return image
}
