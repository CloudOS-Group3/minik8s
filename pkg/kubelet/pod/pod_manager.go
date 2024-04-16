package pod_manager

import (
	"minik8s/pkg/api"
	"minik8s/pkg/kubelet/image"
)

type PodManager struct {
	podByName    map[string]*api.Pod
	ImageManager *image_manager.ImageManager
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
