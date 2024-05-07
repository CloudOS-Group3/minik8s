package pod

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"minik8s/pkg/api"
	"minik8s/pkg/kubelet/container"
	"minik8s/util/log"
)

type PodMetrics struct {
	CpuUsage         float64
	MemoryUsage      float64
	ContainerMetrics []container.ContainerMetrics
}

func CreatePod(pod *api.Pod) bool {

	// create pause container
	pause_container := CreatePauseContainer(pod)
	if pause_container == nil {
		log.Error("Failed to create pause container for pod %s", pod.Metadata.Name)
		return false
	}

	//cAdvisorContainer := CreateCAdvisorContainer(pod)
	//if cAdvisorContainer == nil {
	//	log.Error("Failed to create cAdvisor container for pod %s", pod.Metadata.Name)
	//	return false
	//}

	ctx := namespaces.WithNamespace(context.Background(), pod.Metadata.NameSpace)
	if container.StartContainer(pause_container, ctx) == false {
		return false
	}

	// create other containers
	for _, container_ := range pod.Spec.Containers {
		new_container := container.CreateContainer(container_, pod.Metadata.NameSpace)
		if new_container == nil {
			log.Error("Failed to create container %s", container_.Name)
		}
		if container.StartContainer(new_container, ctx) == false {
			return false
		}
	}

	return true
}

func CreatePauseContainer(pod *api.Pod) containerd.Container {
	config := api.Container{
		Name:            pod.Metadata.Name + "-pause",
		Image:           "registry.cn-hangzhou.aliyuncs.com/google_containers/pause:3.9",
		ImagePullPolicy: api.PullPolicyIfNotPresent,
	}

	return container.CreateContainer(config, pod.Metadata.NameSpace)
}

func CreateCAdvisorContainer(pod *api.Pod) containerd.Container {
	config := api.Container{
		Name:            "cAdvisor",
		Image:           "gcr.io/google-containers/cadvisor:latest",
		ImagePullPolicy: api.PullPolicyIfNotPresent,
		Ports: []api.ContainerPort{
			{
				ContainerPort: 10000,
			},
		},
	}
	return container.CreateContainer(config, pod.Metadata.NameSpace)
}

func GetPodMetrics(pod *api.Pod) (*PodMetrics, error) {
	podMetrics := &PodMetrics{}
	totalCpuUsage := 0.0
	totalMemoryUsage := 0.0

	for _, container_ := range pod.Spec.Containers {
		// fix history bugs
		if pod.Metadata.NameSpace == "" {
			pod.Metadata.NameSpace = "default"
		}
		containerMetrics, err := container.GetContainerMetrics(container_.Name, pod.Metadata.NameSpace)
		if err != nil {
			log.Error("Failed to get metrics for container %s", container_.Name)
			return nil, err
		}
		totalCpuUsage += containerMetrics.CpuUsage
		totalMemoryUsage += containerMetrics.MemoryUsage
		podMetrics.ContainerMetrics = append(podMetrics.ContainerMetrics, *containerMetrics)
	}
	podMetrics.CpuUsage = totalCpuUsage
	podMetrics.MemoryUsage = totalMemoryUsage
	return podMetrics, nil

}

func DeletePod(pod *api.Pod) bool {
	ctx := namespaces.WithNamespace(context.Background(), pod.Metadata.NameSpace)

	// delete containers
	for _, container_ := range pod.Spec.Containers {
		container_to_del := container.GetContainerById(container_.Name, pod.Metadata.NameSpace)
		if container_to_del == nil {
			log.Warn("Container %s not found", container_.Name)
			continue
		}
		if container.RemoveContainer(container_to_del, ctx) == false {
			log.Error("Failed to remove container %s", container_.Name)
			return false
		}
	}

	// delete pause container
	pause_container := container.GetContainerById(pod.Metadata.Name+"-pause", pod.Metadata.NameSpace)
	if pause_container == nil {
		log.Error("Pause container not found")
		return false
	}
	if container.RemoveContainer(pause_container, ctx) == false {
		log.Error("Failed to remove pause container")
		return false
	}

	// delete pod
	return true
}
