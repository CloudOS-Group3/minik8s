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
	"os"
	"os/exec"
	"reflect"
	"strings"
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
func GetPauseName(pod *api.Pod) string {
	return pod.Metadata.NameSpace + "-" + pod.Metadata.Name + "-pause"
}

type ContainerInspect struct {
	State struct {
		Pid int `json:"Pid"`
	} `json:"State"`
	NetworkSettings struct {
		IPAddress string `json:"IPAddress"`
	} `json:"NetworkSettings"`
}

// retrun pause PID
func CreatePauseContainer(pod *api.Pod) (string, error) {
	// get all ports that need to be exposed
	//ports := []api.ContainerPort{}
	//for _, container := range pod.Spec.Containers {
	//	for _, port := range container.Ports {
	//		ports = append(ports, port)
	//	}
	//}
	pause_name := GetPauseName(pod)
	pause_image := "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.9"

	// Use nerdctl to create pause container
	// network: cbr0, which is flannel network
	// the cmd will output the container id
	cmd := exec.Command("nerdctl", "-n", pod.Metadata.NameSpace, "run", "-d", "--name", pause_name, "--network", "cbr0", pause_image)
	log.Info("cmd: %v", cmd)
	containerID, err := cmd.Output()
	if err != nil {
		log.Error("Failed to run nerdctl run: %s", err.Error())
		return "", err
	}
	trimmedContainerID := strings.TrimSpace(string(containerID))
	pod.Status.PauseId = trimmedContainerID
	log.Info("Create pause: %s", trimmedContainerID)

	// get pause container pid & ip
	cmd = exec.Command("nerdctl", "-n", pod.Metadata.NameSpace, "inspect", trimmedContainerID)
	output, err := cmd.Output()
	// unmarshal JSON
	var inspectData []ContainerInspect
	err = json.Unmarshal(output, &inspectData)
	if err != nil {
		log.Error("Failed to parse JSON: %s", err.Error())
	}

	// get first container pid
	pid := inspectData[0].State.Pid
	ip := inspectData[0].NetworkSettings.IPAddress

	// Set pod ip = pause container ip
	pod.Status.PodIP = ip
	log.Info("Pause container ip: %s", ip)

	return fmt.Sprintf("%d", pid), nil
}

func ExecuteCommandInContainer(ctx context.Context, task containerd.Task, command []string) error {
	execID := command[0]
	execProcessSpec := specs.Process{
		Args: command,
		Cwd:  "/",
		Env:  []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"},
	}

	execTask, err := task.Exec(ctx, execID, &execProcessSpec, cio.NewCreator(cio.WithStdio))
	if err != nil {
		return err
	}
	defer execTask.Delete(ctx)

	if err := execTask.Start(ctx); err != nil {
		return err
	}

	statusC, err := execTask.Wait(ctx)
	if err != nil {
		return err
	}

	status := <-statusC
	code, _, err := status.Result()
	if err != nil {
		return err
	}

	if code != 0 {
		return fmt.Errorf("command %s exited with status %d", command, code)
	}

	return nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)

	// Redirect stdout and stderr to os.Stdout and os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Check if we need to provide stdin as a console
	if name == "nerdctl" && len(args) > 0 && args[0] == "exec" {
		cmd.Stdin = os.Stdin
	}

	return cmd.Run()
}

func CreateContainer(config api.Container, namespace string, pause_pid string, hostMount map[string]string, volumes []api.Volume) containerd.Container {
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

	/* Set Linux namespace */
	//allow container to have their own pid namespace
	//opt = append(opt, oci.WithLinuxNamespace(specs.LinuxNamespace{Type: "pid", Path: "/proc/" + pause_pid + "/ns/pid"}))
	opt = append(opt, oci.WithLinuxNamespace(specs.LinuxNamespace{Type: "ipc", Path: "/proc/" + pause_pid + "/ns/ipc"}))
	opt = append(opt, oci.WithLinuxNamespace(specs.LinuxNamespace{Type: "uts", Path: "/proc/" + pause_pid + "/ns/uts"}))
	//opt = append(opt, oci.WithLinuxNamespace(specs.LinuxNamespace{Type: "mount", Path: "/proc/" + pause_pid + "/ns/mount"}))
	opt = append(opt, oci.WithLinuxNamespace(specs.LinuxNamespace{Type: "network", Path: "/proc/" + pause_pid + "/ns/net"}))

	// run as privileged container
	opt = append(opt, oci.WithPrivileged)

	/* Manage volume: mount host path to container */
	log.Info("hostMount: %v", hostMount)
	for _, volume := range config.VolumeMounts {
		permission := "rw"
		if volume.ReadOnly {
			permission = "ro"
		}
		opt = append(opt, oci.WithMounts([]specs.Mount{
			{
				Destination: volume.MountPath,
				Source:      hostMount[volume.Name],
				Type:        "bind",
				Options:     []string{"rbind", permission},
			},
		}))
	}
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

	task, err := container.NewTask(ctx, cio.NewCreator(cio.WithStdio))
	if err != nil {
		log.Error("failed to create container task: %v", err)
		return nil
	}
	defer task.Delete(context.Background())

	if err = task.Start(ctx); err != nil {
		log.Error("failed to start container task: %v", err)
		return nil
	}

	//ExecuteCommandInContainer(ctx, task, []string{"apt", "update"})
	//ExecuteCommandInContainer(ctx, task, []string{"apt", "install", "nfs-common"})
	//ExecuteCommandInContainer(ctx, task, []string{"mkdir", "-p", volumes[0].NFS.Path})
	//ExecuteCommandInContainer(ctx, task, []string{"mount", "-t", "nfs", "192.168.3.7:/nfsroot", volumes[0].NFS.Path})

	if len(volumes) > 0 && volumes[0].NFS.Path != "" {
		if err = runCommand("nerdctl", "exec", "-it", config.Name, "apt", "update"); err != nil {
			log.Error("error apt update")
		}
		if err = runCommand("nerdctl", "exec", "-it", config.Name, "apt", "install", "nfs-common", "-y"); err != nil {
			log.Error("error apt install")
		}
		if err = runCommand("nerdctl", "exec", "-it", config.Name, "mkdir", "-p", volumes[0].NFS.Path); err != nil {
			log.Error("error mkdir")
		}
		if err = runCommand("nerdctl", "exec", "-it", config.Name, "mount", "-t", "nfs", "-o", "nolock", "192.168.3.6:/nfsroot/test", volumes[0].NFS.Path); err != nil {
			log.Error("error mount")
		}
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

func RemoveContainer(name string, namespace string) bool {
	cmd := exec.Command("nerdctl", "-n", namespace, "stop", name)
	_, err := cmd.Output()
	if err != nil {
		log.Error("Failed to stop container %s", err.Error())
	}
	cmd = exec.Command("nerdctl", "-n", namespace, "rm", name)
	_, err = cmd.Output()
	if err != nil {
		log.Error("Failed to remove container %s", err.Error())
		return false
	}
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
	//log.Debug("v: %v", v)
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

func GetContainerIdByName(name string, ctx context.Context) (string, error) {
	client, err := util.CreateClient()
	if err != nil {
		log.Error("Failed to create containerd client: %v", err.Error())
		return "", err
	}
	defer client.Close()

	containers, err := client.Containers(ctx)
	if err != nil {
		return "", err
	}

	for _, container := range containers {
		info, err := container.Info(ctx)
		if err != nil {
			return "", err
		}
		if info.Labels["io.containerd.container.name"] == name {
			return container.ID(), nil
		}
		if info.Labels["nerdctl/name"] == name {
			return container.ID(), nil
		} // pause is created by nerdctl
		log.Debug("label: %v", info.Labels)
	}

	return "", fmt.Errorf("container with name %s not found", name)
}
