package container

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "github.com/containerd/cgroups/stats/v1"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
	"github.com/gogo/protobuf/proto"
	"github.com/opencontainers/runtime-spec/specs-go"
	"minik8s/pkg/api"
	"minik8s/pkg/kubelet/image"
	"minik8s/pkg/util"
	"minik8s/util/log"
	"os/exec"
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
func CreatePauseContainer(pod *api.Pod) containerd.Container {
	// get all ports that need to be exposed
	ports := []api.ContainerPort{}
	for _, container := range pod.Spec.Containers {
		for _, port := range container.Ports {
			ports = append(ports, port)
		}
	}
	config := api.Container{
		Name:            pod.Metadata.Name + "-pause",
		Image:           "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.9",
		ImagePullPolicy: api.PullPolicyIfNotPresent,
		Command:         []string{"pause"},
	}

	client, err := util.CreateClient()
	if err != nil {
		log.Error("Failed to create containerd client: %v", err.Error())
		return nil
	}
	ctx := namespaces.WithNamespace(context.Background(), pod.Metadata.NameSpace)

	// pull image

	image_ := image.PullImage(config.Image, config.ImagePullPolicy, client, pod.Metadata.NameSpace)
	if image_ == nil {
		log.Error("Failed to pull image %s", config.Image)
		return nil
	}

	// create container

	// check if exists
	container_, err := client.LoadContainer(ctx, config.Name)
	if err == nil {
		log.Info("Container %s already exists", config.Name)
		return container_
	}
	container, err := client.NewContainer(
		ctx,
		config.Name,
		containerd.WithNewSnapshot(config.Name+"_"+fmt.Sprintf("%d", time.Now().Unix()), image_),
		containerd.WithNewSpec(oci.WithImageConfig(image_)),
	)
	if err != nil {
		log.Error("Failed to create container %s: %v", config.Name, err.Error())
		return nil
	}
	log.Info("Container %s created", container.ID())
	return container
}

func CreateContainer(config api.Container, namespace string, pause_pid string) containerd.Container {
	client, err := util.CreateClient()
	if err != nil {
		log.Error("Failed to create containerd client: %v", err.Error())
		return nil
	}
	ctx := namespaces.WithNamespace(context.Background(), namespace)

	// pull image

	image_ := image.PullImage(config.Image, config.ImagePullPolicy, client, namespace)
	if image_ == nil {
		log.Error("Failed to pull image %s", config.Image)
		return nil
	}

	// create container

	// check if exists
	container_, err := client.LoadContainer(ctx, config.Name)
	if err == nil {
		log.Info("Container %s already exists", config.Name)
		return container_
	}
	opt := []oci.SpecOpts{oci.WithImageConfig(image_)}
	opt = append(opt, oci.WithLinuxNamespace(specs.LinuxNamespace{Type: "pid", Path: "/proc/" + pause_pid + "/ns/pid"}))
	opt = append(opt, oci.WithLinuxNamespace(specs.LinuxNamespace{Type: "ipc", Path: "/proc/" + pause_pid + "/ns/ipc"}))
	opt = append(opt, oci.WithLinuxNamespace(specs.LinuxNamespace{Type: "uts", Path: "/proc/" + pause_pid + "/ns/uts"}))
	opt = append(opt, oci.WithLinuxNamespace(specs.LinuxNamespace{Type: "network", Path: "/proc/" + pause_pid + "/ns/net"}))
	opt_ := []containerd.NewContainerOpts{
		//containerd.WithImage(image_),
		containerd.WithNewSnapshot(config.Name+"_"+fmt.Sprintf("%d", time.Now().Unix()), image_),
		containerd.WithNewSpec(opt...),
	}
	container, err := client.NewContainer(
		ctx,
		config.Name,
		opt_...,
	)
	if err != nil {
		log.Error("Failed to create container %s: %v", config.Name, err.Error())
		return nil
	}
	log.Info("Container %s created", config.Name)
	return container

}

func GetContainerById(container_id string, namespace string) containerd.Container {
	if namespace == "" {
		namespace = "default"
	}
	client, err := util.CreateClient()
	if err != nil {
		log.Error("Failed to create containerd client: %v", err.Error())
		return nil
	}
	ctx := namespaces.WithNamespace(context.Background(), namespace)
	container_, err := client.LoadContainer(ctx, container_id)
	if err != nil {
		log.Info("Failed to load container %s: %v", container_id, err.Error())
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
			log.Info("Container %s already started", container.ID())
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
		log.Error("Failed to create task for container %s: %v", container.ID(), err.Error())
		return false
	}
	err = task.Start(ctx)
	if err != nil {
		log.Error("Failed to start task for container %s: %v", container.ID(), err.Error())
		return false
	}
	log.Info("Container %s started", container.ID())
	return true
}

func StopContainer(container containerd.Container, ctx context.Context) bool {
	// search for task
	task, err := container.Task(ctx, nil)
	if err != nil {
		log.Error("Failed to get task for container %s: %v", container.ID(), err.Error())
		return false
	}

	// check if already stopped
	status, err := task.Status(ctx)
	if err == nil && status.Status == containerd.Stopped {
		log.Error("Container %s already stopped", container.ID())
		return true
	}

	// kill task
	err = task.Kill(ctx, syscall.SIGTERM)
	if err != nil {
		log.Error("Failed to kill task for container %s: %v", container.ID(), err)
		return false
	}

	// wait for task to exit
	exitStatusC, err := task.Wait(ctx)
	if err != nil {
		log.Error("Failed to wait task for container %s: %v", container.ID(), err.Error())
		return false
	}
	select {
	case <-exitStatusC:
		break
	case <-time.After(30 * time.Second):
		log.Error("Failed to wait task for container %s: timeout", container.ID())
	}
	log.Info("Container %s stopped", container.ID())
	return true
}

func RemoveContainer(container containerd.Container, ctx context.Context) bool {
	// search for task
	task, err := container.Task(ctx, nil)
	if err != nil {
		log.Error("Failed to get task for container %s: %v", container.ID(), err.Error())
		return false
	}
	// check status
	status, err := task.Status(ctx)
	if err != nil {
		log.Error("Failed to get task status for container %s: %v", container.ID(), err.Error())
		return false
	}

	// stop task if running
	if status.Status == containerd.Running {
		err = task.Kill(ctx, syscall.SIGTERM)
		if err != nil {
			log.Error("Failed to kill task for container %s: %v", container.ID(), err.Error())
			return false
		}

		// wait for task to exit
		exitStatusC, err := task.Wait(ctx)
		if err != nil {
			log.Error("Failed to wait task for container %s: %v", container.ID(), err.Error())
			return false
		}
		select {
		case <-exitStatusC:
			break
		case <-time.After(30 * time.Second):
			log.Error("Failed to wait task for container %s: timeout", container.ID())
			return false
		}

	}
	// delete task
	_, err = task.Delete(ctx)
	if err != nil {
		log.Error("Failed to remove task for container %s: %v", container.ID(), err.Error())
		return false
	}

	// remove container
	err = container.Delete(ctx, containerd.WithSnapshotCleanup)
	if err != nil {
		log.Error("Failed to remove container %s: %v", container.ID(), err.Error())
		return false
	}
	log.Info("Container %s removed", container.ID())
	return true
}

func GetContainerMetrics(name string, space string) (*ContainerMetrics, error) {
	client, err := util.CreateClient()
	if err != nil {
		log.Error("Failed to create containerd client: %v", err.Error())
		return nil, err
	}
	defer client.Close()

	ctx := namespaces.WithNamespace(context.Background(), space)
	container, err := client.LoadContainer(ctx, name)
	if err != nil {
		log.Error("Failed to get container %s: %s", name, err.Error())
		return nil, fmt.Errorf("Failed to get container %s", name)
	}

	task, err := container.Task(ctx, nil)
	if err != nil {
		log.Error("Failed to get task for container %s: %v", name, err.Error())
		return nil, err
	}
	status, err := task.Status(ctx)
	if err != nil {
		log.Error("Failed to get task status for container %s: %v", name, err.Error())
		return nil, err
	}
	metrics, err := task.Metrics(ctx)
	if err != nil {
		log.Error("Failed to get metrics for container %s: %v", name, err.Error())
		return nil, err
	}
	// Unmarshal metrics
	// metrics.Data.Value is any type, and is a serialized v1.Metrics, according to fmt.Println(typeurl.UnmarshalAny(metrics.Data))
	// Reference: https://github.com/IPADSIntern-MiniK8s/MiniK8s/blob/master/pkg/kubelet/container/container.go#L209
	v := reflect.New(reflect.TypeOf(v1.Metrics{})).Interface()
	err = proto.Unmarshal(metrics.Data.Value, v.(proto.Message))
	log.Debug("v: %v", v)
	if err != nil {
		log.Error("Failed to unmarshal metrics for container %s: %v", name, err.Error())
		return nil, err
	}

	return &ContainerMetrics{
		CpuUsage:      float64(v.(*v1.Metrics).CPU.Usage.Total),
		MemoryUsage:   float64(v.(*v1.Metrics).Memory.Usage.Usage),
		ProcessStatus: string(status.Status),
	}, nil

}

type ContainerInspect struct {
	State struct {
		Pid int `json:"Pid"`
	} `json:"State"`
}

func GetContainerPid(container containerd.Container, namespace string) string {
	cmd := exec.Command("nerdctl", "-n", namespace, "inspect", container.ID())
	output, err := cmd.Output()
	if err != nil {
		log.Error("Failed to run nerdctl inspect: %s", err.Error())
	}

	// unmarshal JSON
	var inspectData []ContainerInspect
	err = json.Unmarshal(output, &inspectData)
	if err != nil {
		log.Error("Failed to parse JSON: %s", err.Error())
	}

	// get first container pid
	return fmt.Sprintf("%d", inspectData[0].State.Pid)
}
