package controllers

import (
	"encoding/json"
	"io/ioutil"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"net/http"
	"strings"
)

func GetPod(namespace string, name string) (*api.Pod, error) {
	URL := config.GetUrlPrefix() + config.PodURL
	strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
	strings.Replace(URL, config.NamePlaceholder, name, -1)

	res, err := http.Get(URL)
	if err != nil {
		log.Error("err get pod %s:%s", namespace, name)
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	pod := &api.Pod{}
	err = json.Unmarshal(body, &pod)
	if err != nil {
		log.Error("error unmarshal into pod")
		return nil, err
	}

	return pod, nil
}

func GetPodsByNamespace(namespace string) ([]api.Pod, error) {
	URL := config.GetUrlPrefix() + config.PodsURL
	strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)

	res, err := http.Get(URL)
	if err != nil {
		log.Error("err get pods in namespace %s", namespace)
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

func ApplyPod(pod *api.Pod) error {
	byteArr, err := json.Marshal(*pod)
	if err != nil {
		log.Error("error marshal pod")
		return err
	}

	URL := config.GetUrlPrefix() + config.PodsURL
	strings.Replace(URL, config.NamespacePlaceholder, pod.Metadata.NameSpace, -1)

	err = httputil.Post(URL, byteArr)
	if err != nil {
		log.Error("error http post pod")
		return err
	}

	return nil
}
