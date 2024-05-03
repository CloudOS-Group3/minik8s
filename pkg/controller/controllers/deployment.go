package controllers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"minik8s/util/stringutil"
	"net/http"
	"strings"
)

func update() {

	allPods, err := getAllPods()
	if err != nil {
		log.Error("get all pods error")
		return
	}

	deployments, err := getAllDeployments()
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
			if checkLabel(pod, deployment) {
				targetPods = append(targetPods, pod)
				allPodsWithDeployment[pod.Metadata.UUID] = true
			}
		}
		if len(targetPods) < deployment.Spec.Replicas {
			addPod(deployment.Spec.Template, deployment.Metadata, deployment.Spec.Replicas-len(targetPods))
		} else if len(targetPods) > deployment.Spec.Replicas {
			deletePod(targetPods, len(targetPods)-deployment.Spec.Replicas)
		}

		updateDeploymentStatus(targetPods, deployment)
	}

	for _, pod := range allPods {
		if pod.Metadata.Labels["deployment"] == "" {
			continue
		}
		if _, ok := allPodsWithDeployment[pod.Metadata.UUID]; !ok {
			deletePod([]api.Pod{pod}, 1)
		}
	}

}

func getAllPods() ([]api.Pod, error) {

	URL := config.GetUrlPrefix() + config.PodsURL
	strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

	res, err := http.Get(URL)
	if err != nil {
		log.Error("err get all pods")
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	pods := []api.Pod{}
	err = json.Unmarshal(body, &pods)
	if err != nil {
		log.Error("error unmarshal into all pods")
		return nil, err
	}

	return pods, nil
}

func getAllDeployments() ([]api.Deployment, error) {
	URL := config.GetUrlPrefix() + config.DeploymentsURL
	strings.Replace(URL, config.NamespacePlaceholder, "default", -1)

	res, err := http.Get(URL)
	if err != nil {
		log.Error("err get all deployments")
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	deployments := []api.Deployment{}
	err = json.Unmarshal(body, &deployments)
	if err != nil {
		log.Error("error unmarshal into all deployments")
		return nil, err
	}

	return deployments, nil
}

// to return true, just need to match one label
func checkLabel(pod api.Pod, deployment api.Deployment) bool {
	for _, label := range deployment.Spec.Selector.MatchLabels {
		if pod.Metadata.Labels[label] != "" {
			return true
		}
	}
	return false
}

func addPod(template api.PodTemplateSpec, deploymentMetadata api.ObjectMeta, number int) {
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
		}

		_, err = http.Post(URL, config.JsonContent, bytes.NewBuffer(byteArr))
		if err != nil {
			log.Error("error automatically add pod in deployment")
		}
	}
}

func deletePod(targetPods []api.Pod, number int) {

	httpClient := &http.Client{}
	for i := 0; i < number; i++ {
		pod := targetPods[i]
		
		URL := config.GetUrlPrefix() + config.PodURL
		URL = strings.Replace(URL, config.NamespacePlaceholder, pod.Metadata.NameSpace, -1)
		URL = strings.Replace(URL, config.NamePlaceholder, pod.Metadata.Name, -1)
		
		req, err := http.NewRequest("DELETE", URL, nil)
		if err != nil {
			log.Error("create new request error")
		}

		res, err := httpClient.Do(req)
		if err != nil {
			log.Error("error in delete request")
		}
		res.Body.Close()
	}
}

func updateDeploymentStatus(targetPods []api.Pod, deployment api.Deployment) {
	
}
