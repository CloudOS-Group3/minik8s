package controllers

import (
	"encoding/json"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"strings"
	"time"
)

type DeploymentController struct{}

const (
	initialDelay   time.Duration = 10
	updateInterval time.Duration = 30
)

func (DC *DeploymentController) Run() {
	<-time.After(initialDelay)

	for {
		DC.update()
		<-time.After(updateInterval)
	}
}

func (DC *DeploymentController) update() {

	allPods, err := DC.getAllPods()
	if err != nil {
		log.Error("get all pods error")
		return
	}

	deployments, err := DC.getAllDeployments()
	if err != nil {
		log.Error("get all deployments error")
		return
	}

	// UUID -> existence
	// used to delete all pods without deployment
	allPodsWithDeployment := map[string]bool{}

	for _, deployment := range deployments {
		targetPods := []api.Pod{}
		for _, pod := range allPods {
			if DC.checkLabel(pod, deployment) {
				targetPods = append(targetPods, pod)
				allPodsWithDeployment[pod.Metadata.UUID] = true
			}
		}
		if len(targetPods) < deployment.Spec.Replicas {
			DC.addPod(deployment.Spec.Template, deployment.Metadata, deployment.Spec.Replicas-len(targetPods))
		} else if len(targetPods) > deployment.Spec.Replicas {
			DC.deletePod(targetPods, len(targetPods)-deployment.Spec.Replicas)
		}

		DC.updateDeploymentStatus(targetPods, deployment)
	}

	for _, pod := range allPods {
		if pod.Metadata.Labels["deployment"] == "" {
			continue
		}
		if _, ok := allPodsWithDeployment[pod.Metadata.UUID]; !ok {
			DC.deletePod([]api.Pod{pod}, 1)
		}
	}

}

func (DC *DeploymentController) getAllPods() ([]api.Pod, error) {

	URL := config.GetUrlPrefix() + config.PodsURL
	strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

	pods := []api.Pod{}

	err := httputil.Get(URL, pods)
	if err != nil {
		log.Error("error get all pods")
		return nil, err
	}

	return pods, nil
}

func (DC *DeploymentController) getAllDeployments() ([]api.Deployment, error) {
	URL := config.GetUrlPrefix() + config.DeploymentsURL
	strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
	deployments := []api.Deployment{}

	err := httputil.Get(URL, deployments)
	if err != nil {
		log.Error("error get all deployments")
		return nil, err
	}


	return deployments, nil
}

// to return true, just need to match one label
func (DC *DeploymentController) checkLabel(targetPod api.Pod, targetDeployment api.Deployment) bool {
	for _, label := range targetDeployment.Spec.Selector.MatchLabels {
		if targetPod.Metadata.Labels[label] != "" {
			return true
		}
	}
	return false
}

func (DC *DeploymentController) addPod(template api.PodTemplateSpec, deploymentMetadata api.ObjectMeta, number int) {
	log.Info("automatically adding pod in deployment")

	var newPod api.Pod
	newPod.APIVersion = "Pod"
	newPod.Kind = "v1"
	newPod.Metadata = template.Metadata
	newPod.Spec = template.Spec
	newPod.Metadata.Labels["deployment"] = deploymentMetadata.UUID

	basePodName := newPod.Metadata.Name
	baseContainerNames := []string{}
	for _, container := range newPod.Spec.Containers {
		baseContainerNames = append(baseContainerNames, container.Name)
	}

	for i := 0; i < number; i++ {
		newPod.Metadata.Name = basePodName + "-" + stringutil.GenerateRandomString(5)

		for index := range newPod.Spec.Containers {
			newPod.Spec.Containers[index].Name = baseContainerNames[index] + "-" + stringutil.GenerateRandomString(5)
		}

		URL := config.GetUrlPrefix() + config.PodsURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, newPod.Metadata.NameSpace, -1)

		byteArr, err := json.Marshal(newPod)
		if err != nil {
			log.Error("error marshal newpod")
			return
		}

		err = httputil.Post(URL, byteArr)
		if err != nil {
			log.Error("error automatically add pod in deployment")
			return
		}
	}
}

func (DC *DeploymentController) deletePod(targetPods []api.Pod, number int) {

	for i := 0; i < number; i++ {
		pod := targetPods[i]

		URL := config.GetUrlPrefix() + config.PodURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, pod.Metadata.NameSpace, -1)
		URL = strings.Replace(URL, config.NamePlaceholder, pod.Metadata.Name, -1)

		err := httputil.Delete(URL)
		if err != nil {
			log.Error("error deleting pod")
			return
		}
	}
}

func (DC *DeploymentController) updateDeploymentStatus(targetPods []api.Pod, targetDeployment api.Deployment) {

	readyPodNum := 0
	for _, pod := range targetPods {
		if pod.Status.Phase == string(api.PodRunning) {
			readyPodNum++
		}
	}

	targetDeployment.Status.ReadyReplicas = readyPodNum

	URL := config.GetUrlPrefix() + config.DeploymentURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, targetDeployment.Metadata.NameSpace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, targetDeployment.Metadata.Name, -1)

	byteArr, err := json.Marshal(targetDeployment)
	if err != nil {
		log.Error("error marshal targetdeployment")
		return
	}

	err = httputil.Put(URL, byteArr)

	if err != nil {
		log.Error("error creating deployment")
		return
	}

}
