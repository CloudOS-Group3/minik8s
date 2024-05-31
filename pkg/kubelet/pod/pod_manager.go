package pod

import (
	"context"
	"github.com/containerd/containerd/namespaces"
	"minik8s/pkg/api"
	"minik8s/pkg/kubelet/container"
	"minik8s/pkg/kubelet/node"
	"minik8s/util/consul"
	"minik8s/util/log"
)

func CreatePod(pod *api.Pod) bool {

	// create pause container & start it
	pause_container_pid, err := container.CreatePauseContainer(pod)
	if err != nil {
		log.Error("Failed to create pause container for pod %s", pod.Metadata.Name)
		return false
	}

	// for VolumeMounts in hostPath
	hostMount := make(map[string]string)
	for _, volume := range pod.Spec.Volumes {
		hostMount[volume.Name] = volume.HostPath
	}

	// create other containers
	ctx := namespaces.WithNamespace(context.Background(), pod.Metadata.NameSpace)
	for _, container_ := range pod.Spec.Containers {
		new_container := container.CreateContainer(container_, pod.Metadata.NameSpace, pause_container_pid, hostMount)
		if new_container == nil {
			log.Error("Failed to create container %s", container_.Name)
		}
		if container.StartContainer(new_container, ctx) == false {
			return false
		}
	}
	log.Info("add pod %v to check list", pod)
	node.AddPodToCheckList(pod)
	count := 0
	for _, cont := range pod.Spec.Containers {
		portCount := 0
		assigned := false
		for _, port := range cont.Ports {
			if port.Prometheus == true {
				assigned = true
				ID := "user-pod-" + pod.Metadata.NameSpace + "-" + pod.Metadata.Name + "-" + string(rune(count)) + "-" + string(rune(portCount))
				name := ID
				addr := pod.Status.PodIP
				watchPort := port.ContainerPort
				consul.RegisterService(ID, name, addr, int(watchPort))
				portCount++
			}
		}
		if assigned {
			count++
		}
	}

	return true
}

func DeletePod(pod *api.Pod) bool {
	// first delete pod from list
	node.DeletePodInCheckList(pod)
	count := 0
	for _, cont := range pod.Spec.Containers {
		portCount := 0
		assigned := false
		for _, port := range cont.Ports {
			if port.Prometheus == true {
				assigned = true
				ID := "user-pod-" + pod.Metadata.NameSpace + "-" + pod.Metadata.Name + "-" + string(rune(count)) + "-" + string(rune(portCount))
				consul.DeRegisterService(ID)
				portCount++
			}
		}
		if assigned {
			count++
		}
	}

	// delete containers
	for _, container_ := range pod.Spec.Containers {
		if container.RemoveContainer(container_.Name, pod.Metadata.NameSpace) == false {
			log.Warn("Failed to remove container %s", container_.Name)
			continue
		}
	}

	// delete pause container
	if container.RemoveContainer(container.GetPauseName(pod), pod.Metadata.NameSpace) == false {
		log.Warn("Failed to remove pause container %s", container.GetPauseName(pod))
		return false
	}

	return true
}
