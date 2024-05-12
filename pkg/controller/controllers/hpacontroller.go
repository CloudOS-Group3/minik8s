package controllers

import (
	"encoding/json"
	"math"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"strings"
	"time"
)

type HPAController struct{}

const (
	hpaDelay    time.Duration = 10
	hpaInterval time.Duration = 30
)

func (this *HPAController) Run() {
	<-time.After(hpaDelay)

	for {
		this.update()
		<-time.After(hpaInterval)
	}
}

func (this *HPAController) update() {
	allPods, err := this.getAllPods()

	if err != nil {
		log.Error("error getting all pods in hpa")
		return
	}

	allHPAs, err := this.getAllHPAs()

	if err != nil {
		log.Error("error getting all HPAs")
		return
	}

	allPodsWithHPA := map[string]bool{}

	for _, hpa := range allHPAs {
		targetPods := []api.Pod{}
		for _, pod := range allPods {
			if this.checkLabel(pod, hpa) {
				targetPods = append(targetPods, pod)
				allPodsWithHPA[pod.Metadata.Name] = true
			}
		}

		cpuUsage := this.calculatePodsCPUUsage(targetPods)
		memoryUsage := this.calculatePodsCPUUsage(targetPods)

		expectReplicaNum := this.calculateExpectReplicaNum(hpa, cpuUsage, memoryUsage)

		if expectReplicaNum > hpa.Status.CurrentReplicas {
			this.deleteHPAPods(targetPods, expectReplicaNum-hpa.Status.CurrentReplicas)
		}
		if expectReplicaNum < hpa.Status.CurrentReplicas {
			this.addHPAPods(hpa, targetPods[0], hpa.Status.CurrentReplicas-expectReplicaNum)
		}
	}

}

func (this *HPAController) getAllPods() ([]api.Pod, error) {

	URL := config.GetUrlPrefix() + config.PodsURL
	strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

	pods := []api.Pod{}

	err := httputil.Get(URL, pods, "data")
	if err != nil {
		log.Error("error get all pods")
		return nil, err
	}

	return pods, nil
}

func (this *HPAController) getAllHPAs() ([]api.HPA, error) {

	URL := config.GetUrlPrefix() + config.HPAsURL

	strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

	hpas := []api.HPA{}
	err := httputil.Get(URL, hpas, "data")
	if err != nil {
		log.Error("error get all hpas")
		return nil, err
	}
	return hpas, nil
}

// to return true, just need to match one label
func (this *HPAController) checkLabel(targetPod api.Pod, targetHPA api.HPA) bool {
	for _, label := range targetHPA.Spec.Selector.MatchLabels {
		if targetPod.Metadata.Labels[label] != "" {
			return true
		}
	}
	return false
}

func (this *HPAController) calculatePodsCPUUsage(pods []api.Pod) float64 {
	totalUsage := 0.0
	for _, pod := range pods {
		totalUsage += pod.Status.CPUPercentage
	}

	return totalUsage / float64(len(pods))
}

func (this *HPAController) calculatePodsMemoryUsage(pods []api.Pod) float64 {
	totalUsage := 0.0
	for _, pod := range pods {
		totalUsage += pod.Status.MemoryPercentage
	}

	return totalUsage / float64(len(pods))
}

func (this *HPAController) calculateExpectReplicaNum(hpa api.HPA, cpuUsage float64, memoryUsage float64) int {
	cpuRatio := cpuUsage / hpa.Spec.Metrics.CPUPercentage
	memoryRatio := memoryUsage / hpa.Spec.Metrics.MemoryPercentage

	expectNum := int(math.Max(cpuRatio, memoryRatio) * float64(hpa.Status.CurrentReplicas))

	if expectNum > hpa.Spec.MaxReplica {
		expectNum = hpa.Spec.MaxReplica
	}
	if expectNum < hpa.Spec.MinReplica {
		expectNum = hpa.Spec.MinReplica
	}

	return expectNum

}

func (this *HPAController) addHPAPods(hpa api.HPA, template api.Pod, num int) {
	for i := 0; i < num; i++ {
		newPod := template
		newPod.Metadata.Name = template.Metadata.Name + "-" + stringutil.GenerateRandomString(5)
		for index := range newPod.Spec.Containers {
			newPod.Spec.Containers[index].Name = newPod.Spec.Containers[index].Name + "-" + stringutil.GenerateRandomString(5)
		}

		newPod.Metadata.Labels["hpaName"] = hpa.Metadata.Name
		newPod.Metadata.Labels["hpaNamespace"] = hpa.Metadata.NameSpace
		newPod.Metadata.Labels["hpaUUID"] = hpa.Metadata.UUID

		URL := config.GetUrlPrefix() + config.PodsURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

		byteArr, err := json.Marshal(newPod)
		if err != nil {
			log.Error("error marshal new pod: %s", err.Error())
			continue
		}

		err = httputil.Post(URL, byteArr)
		if err != nil {
			log.Debug("error add new pod: %s", err.Error())
			continue
		}

	}
}

func (this *HPAController) deleteHPAPods(targetPods []api.Pod, num int) {
	for i := 0; i < num; i++ {
		pod := targetPods[i]
		URL := config.GetUrlPrefix() + config.PodURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, pod.Metadata.NameSpace, -1)
		URL = strings.Replace(URL, config.NamePlaceholder, pod.Metadata.Name, -1)

		err := httputil.Delete(URL)
		if err != nil {
			log.Error("error delete pod in hpa")
		}
	}
}
