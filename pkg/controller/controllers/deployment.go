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
	initialDelay   time.Duration = 0 * time.Second
	updateInterval time.Duration = 10 * time.Second
)

func (this *DeploymentController) Run() {
	<-time.After(initialDelay)

	for {
		this.update()
		<-time.After(updateInterval)
	}
}

func (this *DeploymentController) update() {

	allPods, err := this.getAllPods()
	log.Debug("before getting all pods")
	if err != nil {
		log.Error("get all pods error")
		return
	}
	log.Debug("all pods are %+v", allPods)
	allDeployments, err := this.getAllDeployments()
	log.Debug("before getting all allDeployments")
	if err != nil {
		log.Error("get all allDeployments error")
		return
	}
	log.Debug("all allDeployments found: %+v", allDeployments)
	// UUID -> existence
	// used to delete all pods without deployment
	allPodsWithDeployment := map[string]bool{}

	for _, deployment := range allDeployments {
		targetPods := []api.Pod{}
		for _, pod := range allPods {
			if this.checkLabel(pod, deployment) {
				targetPods = append(targetPods, pod)
				allPodsWithDeployment[pod.Metadata.UUID] = true
			}
		}
		log.Debug("targetpods: %v", targetPods)
		if len(targetPods) < deployment.Spec.Replicas {
			this.addPod(deployment.Spec.Template, deployment.Metadata, deployment.Spec.Replicas-len(targetPods))
		} else if len(targetPods) > deployment.Spec.Replicas {
			this.deletePod(targetPods, len(targetPods)-deployment.Spec.Replicas)
		}

		this.updateDeploymentStatus(targetPods, deployment)
	}

	for _, pod := range allPods {
		if pod.Metadata.Labels["deployment"] == "" {
			continue
		}
		if _, ok := allPodsWithDeployment[pod.Metadata.UUID]; !ok {
			this.deletePod([]api.Pod{pod}, 1)
		}
	}

}

func (this *DeploymentController) getAllPods() ([]api.Pod, error) {

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

func (this *DeploymentController) getAllDeployments() ([]api.Deployment, error) {
	URL := config.GetUrlPrefix() + config.DeploymentsURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, "default", -1)
	deployments := []api.Deployment{}

	err := httputil.Get(URL, &deployments, "data")
	if err != nil {
		log.Error("error get all deployments")
		return nil, err
	}
	log.Debug("all deployments are: %v", deployments)
	return deployments, nil
}

// to return true, just need to match one label
func (this *DeploymentController) checkLabel(targetPod api.Pod, targetDeployment api.Deployment) bool {
	for key, _ := range targetDeployment.Spec.Selector.MatchLabels {
		if targetPod.Metadata.Labels[key] != "" {
			return true
		}
	}
	return false
}

func (this *DeploymentController) addPod(template api.PodTemplateSpec, deploymentMetadata api.ObjectMeta, number int) {
	log.Info("automatically adding pod in deployment")

	var newPod api.Pod
	newPod.APIVersion = "v1"
	newPod.Kind = "Pod"
	// deep copy struct
	buffer, err := json.Marshal(template)
	if err != nil {
		log.Error("error marshalling pod template")
	}
	err = json.Unmarshal([]byte(buffer), &newPod)
	if err != nil {
		log.Error("error unmarshalling pod template")
	}
	newPod.Metadata.Labels["deployment"] = deploymentMetadata.UUID
	log.Debug("the content of new pod is: %+v", newPod)

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

func (this *DeploymentController) deletePod(targetPods []api.Pod, number int) {

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

func (this *DeploymentController) updateDeploymentStatus(targetPods []api.Pod, targetDeployment api.Deployment) {

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
