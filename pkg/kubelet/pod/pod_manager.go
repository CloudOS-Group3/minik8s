package pod_manager

import (
	"minik8s/pkg/api"
	"minik8s/pkg/kubelet/container"
)

type PodManager struct {
	podByName        map[string]*api.Pod
	ContainerManager *container.ContainerManager
}

func (pm *PodManager) CreatePod(pod *api.Pod) {
	for _, container_ := range pod.Spec.Containers {
		pm.ContainerManager.CreateContainer(container_)
	}
}

//func (pm *PodManager) GetPodByName(name string) *api.Pod {
//	return pm.podByName[name]
//}
//
//func (pm *PodManager) AddPod(pod *api.Pod) {
//	pm.podByName[pod.Metadata.Name] = pod
//}
//
//func (pm *PodManager) DeletePodByName(name string) {
//	delete(pm.podByName, name)
//}
