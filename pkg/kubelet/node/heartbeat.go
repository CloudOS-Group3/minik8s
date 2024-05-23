package node

import (
	"encoding/json"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/pkg/kubelet/container"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"strings"
	"time"
)

type NodeInfo struct {
	Node *api.Node
}

var Heartbeat *NodeInfo = nil

func init() {
	NewNode := &api.Node{}
	NewNode.APIVersion = "v1"
	NewNode.Kind = "Node"
	NewNode.Metadata.Name = "node1" // TODO: config
	NewNode.Status.Pods = make([]api.Pod, 0)
	NewNode.Status.PodsNumber = 0
	Heartbeat = &NodeInfo{NewNode}
}

func AddPodToCheckList(pod *api.Pod) {
	dupe := false
	for index, PodInList := range Heartbeat.Node.Status.Pods {
		if PodInList.Metadata.Name == pod.Metadata.Name && PodInList.Metadata.NameSpace == pod.Metadata.NameSpace {
			Heartbeat.Node.Status.Pods[index] = *pod
			dupe = true
		}
	}
	if dupe {
		return
	}
	Heartbeat.Node.Status.Pods = append(Heartbeat.Node.Status.Pods, *pod)
	Heartbeat.Node.Status.PodsNumber++
}

func DeletePodInCheckList(pod *api.Pod) {
	for index, PodInList := range Heartbeat.Node.Status.Pods {
		if PodInList.Metadata.Name == pod.Metadata.Name && PodInList.Metadata.NameSpace == pod.Metadata.NameSpace {
			NewList := append(Heartbeat.Node.Status.Pods[:index], Heartbeat.Node.Status.Pods[index+1:]...)
			Heartbeat.Node.Status.Pods = NewList
			Heartbeat.Node.Status.PodsNumber--
			break
		}
	}
}

func DoHeartBeat() {
	for {
		log.Info("Start HeartBeat, Pod number: %d", Heartbeat.Node.Status.PodsNumber)
		for index, PodInList := range Heartbeat.Node.Status.Pods {
			Metrics, err := GetPodMetrics(&PodInList)
			if err != nil {
				panic(err)
			}
			PodInList.Status.Metrics = *Metrics
			PodInList.Status.CPUPercentage = Metrics.CpuUsage
			PodInList.Status.MemoryPercentage = Metrics.MemoryUsage
			URL := config.GetUrlPrefix() + config.PodURL
			URL = strings.Replace(URL, config.NamespacePlaceholder, PodInList.Metadata.NameSpace, -1)
			URL = strings.Replace(URL, config.NamePlaceholder, PodInList.Metadata.Name, -1)
			byteArr, err := json.Marshal(PodInList)
			log.Info("HeartBeat, Pod: %s", string(byteArr))
			if err != nil {
				panic(err)
			}
			//log.Info("HeartBeat, Pod: %s", string(byteArr))
			err = httputil.Put(URL, byteArr)
			if err != nil {
				panic(err)
			}
			Heartbeat.Node.Status.Pods[index] = PodInList
		}
		Heartbeat.Node.Status.Condition.LastHeartbeatTime = time.Now()
		Heartbeat.Node.Status.Condition.Status = "Ready"
		URL := config.GetUrlPrefix() + config.NodeURL
		URL = strings.Replace(URL, config.NamePlaceholder, Heartbeat.Node.Metadata.Name, -1)
		byteArr, err := json.Marshal(Heartbeat.Node)
		//log.Info("HeartBeat, Node: %s", string(byteArr))
		if err != nil {
			panic(err)
		}
		err = httputil.Put(URL, byteArr)
		if err != nil {
			panic(err)
		}
		log.Info("HeartBeat done")
		time.Sleep(30 * time.Second)
	}
}

func GetPodMetrics(pod *api.Pod) (*api.PodMetrics, error) {
	podMetrics := &api.PodMetrics{}
	totalCpuUsage := 0.0
	totalMemoryUsage := 0.0

	for _, container_ := range pod.Spec.Containers {
		// fix history bugs
		if pod.Metadata.NameSpace == "" {
			pod.Metadata.NameSpace = "default"
		}
		containerMetrics, err := container.GetContainerMetrics(container_.Name, pod.Metadata.NameSpace)
		if err != nil {
			log.Info("Failed to get metrics for container %s", container_.Name)
			continue
		}
		totalCpuUsage += containerMetrics.CpuUsage
		totalMemoryUsage += containerMetrics.MemoryUsage
		metrics := &api.ContainerMetrics{
			CpuUsage:      containerMetrics.CpuUsage,
			MemoryUsage:   containerMetrics.MemoryUsage,
			ProcessStatus: containerMetrics.ProcessStatus,
		}
		podMetrics.ContainerMetrics = append(podMetrics.ContainerMetrics, *metrics)
	}
	podMetrics.CpuUsage = totalCpuUsage
	podMetrics.MemoryUsage = totalMemoryUsage
	return podMetrics, nil

}
