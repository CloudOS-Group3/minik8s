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

var AllHPAWaitTime map[string]float64
var AllHPAReady map[string]bool

const (
	hpaInterval time.Duration = 5 * time.Second
)

func (this *HPAController) Run() {
	AllHPAWaitTime = make(map[string]float64)
	AllHPAReady = make(map[string]bool)
	for {
		allHPAs, err := this.getAllHPAs()
		if err != nil {
			log.Error("Failed to get HPA resources")
		}
		for _, hpa := range allHPAs {
			AllHPAWaitTime[hpa.Metadata.Name] += hpaInterval.Seconds()
			if AllHPAWaitTime[hpa.Metadata.Name] >= hpa.Spec.AdjustInterval {
				AllHPAReady[hpa.Metadata.Name] = true
				AllHPAWaitTime[hpa.Metadata.Name] -= hpa.Spec.AdjustInterval
			}
		}
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
	log.Debug("all pods are: %v", allPods)

	allHPAs, err := this.getAllHPAs()

	if err != nil {
		log.Error("error getting all HPAs")
		return
	}

	log.Debug("all hpas are: %v", allHPAs)

	for _, hpa := range allHPAs {
		if !AllHPAReady[hpa.Metadata.Name] {
			continue
		}
		AllHPAReady[hpa.Metadata.Name] = false
		targetPods := []api.Pod{}
		for _, pod := range allPods {
			if this.checkLabel(pod, hpa) {
				targetPods = append(targetPods, pod)
			}
		}

		cpuUsage := this.calculatePodsCPUUsage(targetPods)
		memoryUsage := this.calculatePodsMemoryUsage(targetPods)

		expectReplicaNum := this.calculateExpectReplicaNum(hpa, cpuUsage, memoryUsage)

		if expectReplicaNum < len(targetPods) {
			this.deleteHPAPods(targetPods, len(targetPods)-expectReplicaNum)
		}
		if expectReplicaNum > len(targetPods) {
			this.addHPAPods(hpa, hpa.Spec.Template, expectReplicaNum-len(targetPods))
		}

		hpa.Status.CurrentReplicas = expectReplicaNum
		URL := config.GetUrlPrefix() + config.HPAURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
		URL = strings.Replace(URL, config.NamePlaceholder, hpa.Metadata.Name, -1)
		bytes, err := json.Marshal(hpa)
		if err != nil {
			log.Error("error marshalling HPA")
			continue
		}
		err = httputil.Put(URL, bytes)
		if err != nil {
			log.Error("error putting HPA")
			continue
		}
	}

}

func (this *HPAController) getAllPods() ([]api.Pod, error) {

	URL := config.GetUrlPrefix() + config.PodsURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

	pods := []api.Pod{}

	err := httputil.Get(URL, &pods, "data")
	if err != nil {
		log.Error("error get all pods")
		return nil, err
	}

	return pods, nil
}

func (this *HPAController) getAllHPAs() ([]api.HPA, error) {

	URL := config.GetUrlPrefix() + config.HPAsURL

	URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

	hpas := []api.HPA{}
	err := httputil.Get(URL, &hpas, "data")
	if err != nil {
		log.Error("error get all hpas")
		return nil, err
	}
	return hpas, nil
}

// to return true, just need to match one label
func (this *HPAController) checkLabel(targetPod api.Pod, targetHPA api.HPA) bool {
	for key, value := range targetHPA.Spec.Selector.MatchLabels {
		if targetPod.Metadata.Labels[key] == value {
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

	log.Debug("CPUPercentage: %v, MemoryPercentage: %v", hpa.Spec.Metrics.CPUPercentage, hpa.Spec.Metrics.MemoryPercentage)
	log.Debug("cpuUsage: %v, memoryUsage: %v", cpuRatio, memoryRatio)

	expectNum := int(math.Max(cpuRatio, memoryRatio) * float64(hpa.Status.CurrentReplicas))

	if expectNum > hpa.Spec.MaxReplica {
		expectNum = hpa.Spec.MaxReplica
	}
	if expectNum < hpa.Spec.MinReplica {
		expectNum = hpa.Spec.MinReplica
	}

	log.Debug("expect replica num is %d", expectNum)
	return expectNum

}

func (this *HPAController) addHPAPods(hpa api.HPA, template api.PodTemplateSpec, num int) {
	log.Debug("adding new hpa pods")
	log.Debug("num is %d", num)
	for i := 0; i < num; i++ {
		var newPod api.Pod
		byteArr, err := json.Marshal(template)
		if err != nil {
			log.Error("error marshalling template")
			continue
		}
		err = json.Unmarshal(byteArr, &newPod)
		if err != nil {
			log.Error("error unmarshalling template")
			continue
		}

		newPod.Metadata.Name = template.Metadata.Name + "-" + stringutil.GenerateRandomString(5)

		for index := range newPod.Spec.Containers {
			newPod.Spec.Containers[index].Name = newPod.Spec.Containers[index].Name + "-" + stringutil.GenerateRandomString(5)
		}

		newPod.Metadata.Labels["hpaUUID"] = hpa.Metadata.UUID

		log.Debug("the pod to be add is %+v", newPod)
		URL := config.GetUrlPrefix() + config.PodsURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

		byteArr, err = json.Marshal(newPod)
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
