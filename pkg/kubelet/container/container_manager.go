package container

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"log"
	"minik8s/pkg/api"
	"minik8s/pkg/kubelet/image"
	"minik8s/pkg/util"
	"time"
)

type ContainerManager struct {
	im *image.ImageManager
}

func (cm *ContainerManager) CreateContainer(config api.Container) containerd.Container {
	client, err := util.CreateClient()
	if err != nil {
		log.Printf("Failed to create containerd client: %v", err)
		return nil
	}
	ctx := namespaces.WithNamespace(context.Background(), "default")

	// pull image
	image_ := cm.im.PullImage(config.Image, config.ImagePullPolicy, client)
	if image_ == nil {
		log.Printf("Failed to pull image %s", config.Image)
		return nil
	}

	// create container
	container, err := client.NewContainer(
		ctx,
		config.Name,
		containerd.WithImage(image_),
		containerd.WithNewSnapshot(config.Name+"_"+fmt.Sprintf("%d", time.Now().Unix()), image_),
		containerd.WithNewSpec(oci.WithImageConfig(image_)),
	)
	if err != nil {
		log.Printf("Failed to create container %s: %v", config.Name, err.Error())
		return nil
	}
	return container

}
