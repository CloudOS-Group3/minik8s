package controllers

import (
	"encoding/json"
	"io/ioutil"
	"minik8s/pkg/api"
	"minik8s/pkg/config"
	"minik8s/util/log"
	"net/http"
	"strings"
)

func GetService(namespace string, name string) (*api.Service, error) {
	URL := config.GetUrlPrefix() + config.ServiceURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, name, -1)
	res, err := http.Get(URL)
	if err != nil {
		log.Error("err get service %s:%s", namespace, name)
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	svc := &api.Service{}
	err = json.Unmarshal(body, &svc)
	if err != nil {
		log.Error("error unmarshal into all deployments")
		return nil, err
	}

	return svc, nil
}

func GetAllServices() ([]api.Service, error) {
	URL := config.GetUrlPrefix() + config.ServicesURL

	res, err := http.Get(URL)
	if err != nil {
		log.Error("err get all services")
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	services := []api.Service{}
	err = json.Unmarshal(body, &services)
	if err != nil {
		log.Error("error unmarshal into all services")
		return nil, err
	}

	return services, nil
}

func GetServicesByNamespace(namespace string) ([]api.Service, error) {
	URL := config.GetUrlPrefix() + config.ServicesURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)

	res, err := http.Get(URL)
	if err != nil {
		log.Error("err get all services")
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	services := []api.Service{}
	err = json.Unmarshal(body, &services)
	if err != nil {
		log.Error("error unmarshal into all services")
		return nil, err
	}

	return services, nil
}

func UpdateService(service *api.Service) error {
	serviceByteArray, err := json.Marshal(service)
	if err != nil {
		return err
	}

	URL := config.GetUrlPrefix() + config.ServiceURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, service.Metadata.NameSpace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, service.Metadata.Name, -1)

	req, err := http.NewRequest(http.MethodPost, URL, strings.NewReader(string(serviceByteArray)))
	if err != nil {
		log.Error("err add service %s:%s", service.Metadata.NameSpace, service.Metadata.Name)
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("err add service %s:%s", service.Metadata.NameSpace, service.Metadata.Name)
		return err
	}

	defer res.Body.Close()

	return nil
}

func DeleteService(namespace string, name string) error {
	URL := config.GetUrlPrefix() + config.ServiceURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, name, -1)

	req, err := http.NewRequest(http.MethodDelete, URL, nil)
	if err != nil {
		log.Error("err delete service %s:%s", namespace, name)
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error("err delete service %s:%s", namespace, name)
		return err
	}

	defer res.Body.Close()

	return nil
}
