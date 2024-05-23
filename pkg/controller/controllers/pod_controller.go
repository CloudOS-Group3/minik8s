package controllers

import (
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"strings"
)

func GetPod(namespace string, name string) (*api.Pod, error) {
	pod := &api.Pod{}
	URL := config.GetUrlPrefix() + config.PodURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, name, -1)

	err := httputil.Get(URL, pod, "data")

	if err != nil {
		log.Error("error get pod: %s", err.Error())
		return nil, err
	}
	return pod, err
}
