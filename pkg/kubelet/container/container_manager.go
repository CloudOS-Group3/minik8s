package container

import (
	"context"
	"fmt"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
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

func NewContainerManager() *ContainerManager {
	im := &image.ImageManager{}
	return &ContainerManager{
		im: im,
	}
}

func (cm *ContainerManager) CreateContainer(config api.Container, namespace string) containerd.Container {
	client, err := util.CreateClient()
	if err != nil {
		log.Printf("Failed to create containerd client: %v", err.Error())
		return nil
	}
	ctx := namespaces.WithNamespace(context.Background(), namespace)

	// pull image

	image_ := cm.im.PullImage(config.Image, config.ImagePullPolicy, client, namespace)
	if image_ == nil {
		log.Printf("Failed to pull image %s", config.Image)
		return nil
	}

	// create container

	// check if exists
	container_, err := client.LoadContainer(ctx, config.Name)
	if err == nil {
		log.Printf("Container %s already exists", config.Name)
		return container_
	}
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

func (*ContainerManager) StartContainerById(container_id string, namespace string) bool {
	client, err := util.CreateClient()
	if err != nil {
		log.Printf("Failed to create containerd client: %v", err.Error())
		return false
	}
	ctx := namespaces.WithNamespace(context.Background(), namespace)
	container_, err := client.LoadContainer(ctx, container_id)
	if err != nil {
		log.Printf("Failed to load container %s: %v", container_id, err.Error())
		return false
	}

	// check if already started
	_, err = container_.Task(ctx, nil)
	if err == nil {
		log.Printf("Container %s already started", container_id)
		return true
	}

	task, err := container_.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		log.Printf("Failed to create task for container %s: %v", container_id, err.Error())
		return false
	}
	err = task.Start(ctx)
	if err != nil {
		log.Printf("Failed to start task for container %s: %v", container_id, err)
		return false
	}
	log.Printf("Container %s: %s started", namespace, container_id)
	return true
}

func (*ContainerManager) StartContainer(container containerd.Container) bool {
	ctx := context.Background()
	task, err := container.NewTask(ctx, nil)
	if err != nil {
		log.Printf("Failed to create task for container %s: %v", container.ID(), err.Error())
		return false
	}
	err = task.Start(ctx)
	if err != nil {
		log.Printf("Failed to start task for container %s: %v", container.ID(), err)
		return false
	}
	log.Printf("Container %s started", container.ID())
	return true
}
