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
	"syscall"
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

func (*ContainerManager) GetContainerById(container_id string, namespace string) containerd.Container {
	client, err := util.CreateClient()
	if err != nil {
		log.Printf("Failed to create containerd client: %v", err.Error())
		return nil
	}
	ctx := namespaces.WithNamespace(context.Background(), namespace)
	container_, err := client.LoadContainer(ctx, container_id)
	if err != nil {
		log.Printf("Failed to load container %s: %v", container_id, err.Error())
		return nil
	}
	return container_
}

func (*ContainerManager) StartContainer(container containerd.Container, ctx context.Context) bool {
	// check if already started
	tasks, _ := container.Task(ctx, nil)
	if tasks != nil {
		status, _ := tasks.Status(ctx)
		//log.Printf("Container %s status: %v", container.ID(), status.Status)
		if status.Status == containerd.Running {
			log.Printf("Container %s already started", container.ID())
			return true
		}
		tasks.Delete(ctx, containerd.WithProcessKill)

		//// why it can't start a stopped task?
		//err :=  tasks.Start(ctx)
		//if err != nil {
		//	log.Printf("Failed to start task for container %s: %v", container.ID(), err.Error())
		//	return false
		//}
		//return true
	}

	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		log.Printf("Failed to create task for container %s: %v", container.ID(), err.Error())
		return false
	}
	err = task.Start(ctx)
	if err != nil {
		log.Printf("Failed to start task for container %s: %v", container.ID(), err.Error())
		return false
	}
	log.Printf("Container %s started", container.ID())
	return true
}

func (*ContainerManager) StopContainer(container containerd.Container, ctx context.Context) bool {
	// search for task
	task, err := container.Task(ctx, nil)
	if err != nil {
		log.Printf("Failed to get task for container %s: %v", container.ID(), err.Error())
		return false
	}

	// check if already stopped
	status, err := task.Status(ctx)
	if err == nil && status.Status == containerd.Stopped {
		log.Printf("Container %s already stopped", container.ID())
		return true
	}

	// kill task
	err = task.Kill(ctx, syscall.SIGTERM)
	if err != nil {
		log.Printf("Failed to kill task for container %s: %v", container.ID(), err)
		return false
	}

	// wait for task to exit
	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		log.Printf("Failed to wait task for container %s: %v", container.ID(), err.Error())
		return false
	}
	select {
	case <-exitStatusC:
		break
	case <-time.After(30 * time.Second):
		log.Printf("Failed to wait task for container %s: timeout", container.ID())
	}
	log.Printf("Container %s stopped", container.ID())
	return true
}

func (*ContainerManager) RemoveContainer(container containerd.Container, ctx context.Context) bool {
	// search for task
	task, err := container.Task(ctx, nil)
	if err != nil {
		log.Printf("Failed to get task for container %s: %v", container.ID(), err.Error())
		return false
	}
	// check status
	status, err := task.Status(ctx)
	if err != nil {
		log.Printf("Failed to get task status for container %s: %v", container.ID(), err.Error())
		return false
	}

	// stop task if running
	if status.Status == containerd.Running {
		err = task.Kill(ctx, syscall.SIGTERM)
		if err != nil {
			log.Printf("Failed to kill task for container %s: %v", container.ID(), err.Error())
			return false
		}

		// wait for task to exit
		exitStatusC, err := task.Wait(ctx)
		if err != nil {
			log.Printf("Failed to wait task for container %s: %v", container.ID(), err.Error())
			return false
		}
		select {
		case <-exitStatusC:
			break
		case <-time.After(30 * time.Second):
			log.Printf("Failed to wait task for container %s: timeout", container.ID())
			return false
		}

	}
	// delete task
	_, err = task.Delete(ctx)
	if err != nil {
		log.Printf("Failed to remove task for container %s: %v", container.ID(), err.Error())
		return false
	}

	// remove container
	err = container.Delete(ctx, containerd.WithSnapshotCleanup)
	if err != nil {
		log.Printf("Failed to remove container %s: %v", container.ID(), err.Error())
		return false
	}
	log.Printf("Container %s removed", container.ID())
	return true
}
