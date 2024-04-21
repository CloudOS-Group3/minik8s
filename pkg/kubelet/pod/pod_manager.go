package pod

import (
	"context"
	"github.com/containerd/containerd"
	"github.com/containerd/containerd/namespaces"
	"log"
	"minik8s/pkg/api"
	"minik8s/pkg/kubelet/container"
)

type PodManager struct {
	podByName        map[string]*api.Pod
	ContainerManager *container.ContainerManager
}

func NewPodManager() *PodManager {
	cm := container.NewContainerManager()

	podManager := &PodManager{
		ContainerManager: cm,
		podByName:        make(map[string]*api.Pod),
	}

	return podManager
}

func (pm *PodManager) CreatePod(pod *api.Pod) bool {

	// create pause container
	pause_container := pm.CreatePauseContainer(pod)
	if pause_container == nil {
		log.Printf("Failed to create pause container for pod %s", pod.Metadata.Name)
		return false
	}
	ctx := namespaces.WithNamespace(context.Background(), pod.Metadata.NameSpace)
	if pm.ContainerManager.StartContainer(pause_container, ctx) == false {
		return false
	}

	// create other containers
	for _, container_ := range pod.Spec.Containers {
		new_container := pm.ContainerManager.CreateContainer(container_, pod.Metadata.NameSpace)
		if new_container == nil {
			log.Printf("Failed to create container %s", container_.Name)
		}
		if pm.ContainerManager.StartContainer(new_container, ctx) == false {
			return false
		}
	}

	// add pod to pod manager
	pm.AddPod(pod)
	return true
}

func (pm *PodManager) CreatePauseContainer(pod *api.Pod) containerd.Container {
	config := api.Container{
		Name:            pod.Metadata.Name + "-pause",
		Image:           "registry.aliyuncs.com/google_containers/pause:3.2",
		ImagePullPolicy: api.PullPolicyIfNotPresent,
	}

	return pm.ContainerManager.CreateContainer(config, pod.Metadata.NameSpace)
}

func (pm *PodManager) ShowPodInfo(name string) {
	pod := pm.GetPodByName(name)
	if pod == nil {
		log.Printf("Pod %s not found", name)
		return
	}
	log.Printf("Pod %s info:", name)
	log.Printf("Namespace: %s", pod.Metadata.NameSpace)
	//log.Printf("UID: %s", pod.Metadata.UID)
	//log.Printf("ResourceVersion: %s", pod.Metadata.ResourceVersion)
	//log.Printf("NodeName: %s", pod.Spec.NodeName)
	//log.Printf("NodeSelector: %v", pod.Spec.NodeSelector)
	log.Printf("Containers:")
	for _, container_ := range pod.Spec.Containers {
		log.Printf("  Name: %s", container_.Name)
		log.Printf("  Image: %s", container_.Image)
		log.Printf("  ImagePullPolicy: %s", container_.ImagePullPolicy)
		//log.Printf("  Ports: %v", container_.Ports)
		//log.Printf("  Args: %v", container_.Args)
		//log.Printf("  Command: %v", container_.Command)
		//log.Printf("  Env: %v", container_.Env)
		//log.Printf("  Resources: %v", container_.Resources)
		//log.Printf("  VolumeMounts: %v", container_.VolumeMounts)
	}

}

func (pm *PodManager) GetPodByName(name string) *api.Pod {
	return pm.podByName[name]
}

func (pm *PodManager) AddPod(pod *api.Pod) {
	pm.podByName[pod.Metadata.Name] = pod
}

func (pm *PodManager) DeletePodByName(name string) bool {
	pod := pm.GetPodByName(name)
	ctx := namespaces.WithNamespace(context.Background(), pod.Metadata.NameSpace)

	// delete containers
	for _, container_ := range pod.Spec.Containers {
		container_to_del := pm.ContainerManager.GetContainerById(container_.Name, pod.Metadata.NameSpace)
		if container_to_del == nil {
			log.Printf("Container %s not found", container_.Name)
			continue
		}
		if pm.ContainerManager.RemoveContainer(container_to_del, ctx) == false {
			log.Printf("Failed to remove container %s", container_.Name)
			return false
		}
	}

	// delete pause container
	pause_container := pm.ContainerManager.GetContainerById(pod.Metadata.Name+"-pause", pod.Metadata.NameSpace)
	if pause_container == nil {
		log.Printf("Pause container not found")
		return false
	}
	if pm.ContainerManager.RemoveContainer(pause_container, ctx) == false {
		log.Printf("Failed to remove pause container")
		return false
	}

	// delete pod
	delete(pm.podByName, name)
	return true
}
