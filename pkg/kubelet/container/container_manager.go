package container

import (
	"context"
	"fmt"
	v1 "github.com/containerd/cgroups/stats/v1"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/gogo/protobuf/proto"
	"log"
	"minik8s/pkg/api"
	"minik8s/pkg/kubelet/image"
	"minik8s/pkg/util"
	"reflect"
	"syscall"
	"time"
)

type ContainerMetrics struct {
	CpuUsage      float64 `protobuf:"fixed64,1"`
	MemoryUsage   float64 `protobuf:"fixed64,2"`
	ProcessStatus string  `protobuf:"bytes,3"`
	//Running ProcessStatus = "running"
	//Created ProcessStatus = "created"
	//Stopped ProcessStatus = "stopped"
	//Paused ProcessStatus = "paused"
	//Pausing ProcessStatus = "pausing"
	//Unknown ProcessStatus = "unknown"
}

func (c *ContainerMetrics) Reset() {
	c.CpuUsage = 0
	c.MemoryUsage = 0
	c.ProcessStatus = ""
}

func (c *ContainerMetrics) String() string {
	return fmt.Sprintf("ContainerMetrics{CpuUsage: %f, MemoryUsage: %f, ProcessStatus: %s}", c.CpuUsage, c.MemoryUsage, c.ProcessStatus)
}

func (c *ContainerMetrics) ProtoMessage() {

}

func CreateContainer(config api.Container, namespace string) containerd.Container {
	client, err := util.CreateClient()
	if err != nil {
		log.Printf("Failed to create containerd client: %v", err.Error())
		return nil
	}
	ctx := namespaces.WithNamespace(context.Background(), namespace)

	// pull image

	image_ := image.PullImage(config.Image, config.ImagePullPolicy, client, namespace)
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

func GetContainerById(container_id string, namespace string) containerd.Container {
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

func StartContainer(container containerd.Container, ctx context.Context) bool {
	// check if already started
	tasks, _ := container.Task(ctx, nil)
	if tasks != nil {
		status, _ := tasks.Status(ctx)
		//log.Printf("Container %s status: %v", container.ID(), status.Status)
		if status.Status == containerd.Running {
			log.Printf("Container %s already started", container.ID())
			return true
		}
		_, err := tasks.Delete(ctx, containerd.WithProcessKill)
		if err != nil {
			return false
		}

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

func StopContainer(container containerd.Container, ctx context.Context) bool {
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

func RemoveContainer(container containerd.Container, ctx context.Context) bool {
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

func GetContainerMetrics(name string, space string) (*ContainerMetrics, error) {
	client, err := util.CreateClient()
	if err != nil {
		log.Printf("Failed to create containerd client: %v", err.Error())
		return nil, err
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), space)
	container, err := client.LoadContainer(ctx, name)
	if err != nil {
		log.Printf("Failed to get container %s", name)
		return nil, fmt.Errorf("Failed to get container %s", name)
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		log.Printf("Failed to get task for container %s: %v", name, err.Error())
		return nil, err
	}
	status, err := task.Status(ctx)
	if err != nil {
		log.Printf("Failed to get task status for container %s: %v", name, err.Error())
		return nil, err
	}
	metrics, err := task.Metrics(ctx)
	if err != nil {
		log.Printf("Failed to get metrics for container %s: %v", name, err.Error())
		return nil, err
	}
	// Unmarshal metrics
	// metrics.Data.Value is any type, and is a serialized v1.Metrics, according to fmt.Println(typeurl.UnmarshalAny(metrics.Data))
	// Reference: https://github.com/IPADSIntern-MiniK8s/MiniK8s/blob/master/pkg/kubelet/container/container.go#L209
	v := reflect.New(reflect.TypeOf(v1.Metrics{})).Interface()
	err = proto.Unmarshal(metrics.Data.Value, v.(proto.Message))
	log.Printf("v: %v", v)
	if err != nil {
		log.Printf("Failed to unmarshal metrics for container %s: %v", name, err.Error())
		return nil, err
	}

	return &ContainerMetrics{
		CpuUsage:      float64(v.(*v1.Metrics).CPU.Usage.Total),
		MemoryUsage:   float64(v.(*v1.Metrics).Memory.Usage.Usage),
		ProcessStatus: string(status.Status),
	}, nil

}
