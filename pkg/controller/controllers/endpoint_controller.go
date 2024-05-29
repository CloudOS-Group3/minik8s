package controllers

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"io/ioutil"
	"minik8s/pkg/api"
	msg_type "minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/kafka"
	"minik8s/pkg/util"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type EndPointController struct {
	subscriber *kafka.Subscriber
	ready      chan bool
	done       chan bool
}

func (e *EndPointController) Setup(session sarama.ConsumerGroupSession) error {
	close(e.ready)
	return nil
}

func (e *EndPointController) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (e *EndPointController) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if msg.Topic == msg_type.PodTopic {
			session.MarkMessage(msg, "")
			podMsg := &msg_type.PodMsg{}
			err := json.Unmarshal(msg.Value, podMsg)
			if err != nil {
				log.Error("unmarshal pod error")
				continue
			}
			switch podMsg.Opt {
			case msg_type.Update:
				// discard pod without pod ip
				if podMsg.NewPod.Status.PodIP == "" {
					break
				}
				// 1. the pod first created by kubelet and the pod ip is assigned
				// 2. the pod's label is changed
				if podMsg.OldPod.Status.PodIP == "" {
					OnPodAdd(&podMsg.NewPod)
				} else if !util.IsLabelEqual(podMsg.NewPod.Metadata.Labels, podMsg.OldPod.Metadata.Labels) {
					OnPodUpdate(&podMsg.NewPod, podMsg.OldPod.Metadata.Labels)
				}
				break
			case msg_type.Delete:
				OnPodDelete(&podMsg.OldPod)
				break
			case msg_type.Add:
				// discard pod without pod ip
				if podMsg.NewPod.Status.PodIP == "" {
					break
				}
				OnPodUpdate(&podMsg.NewPod, nil)
				break
			}
		} else if msg.Topic == msg_type.ServiceTopic {
			session.MarkMessage(msg, "")
			serviceMsg := &msg_type.ServiceMsg{}
			err := json.Unmarshal(msg.Value, serviceMsg)
			if err != nil {
				log.Error("unmarshal service error")
				continue
			}

			log.Info("service message, %v", serviceMsg)
			switch serviceMsg.Opt {
			case msg_type.Update:
				if !util.IsLabelEqual(serviceMsg.NewService.Spec.Selector, serviceMsg.OldService.Spec.Selector) {
					OnServiceUpdate(&serviceMsg.NewService, serviceMsg.OldService.Spec.Selector)
				}
				break
			case msg_type.Delete:
				OnServiceDelete(&serviceMsg.OldService)
				break
			case msg_type.Add:
				OnServiceUpdate(&serviceMsg.NewService, nil)
				break
			}
		}
	}
	return nil
}

func NewEndPointController() *EndPointController {
	group := "endpoint-controller"
	return &EndPointController{
		ready:      make(chan bool),
		done:       make(chan bool),
		subscriber: kafka.NewSubscriber(group),
	}
}
func OnPodAdd(pod *api.Pod) {
	log.Info("OnPodAdd")
	// deal with empty label
	if pod.Metadata.Labels == nil || len(pod.Metadata.Labels) == 0 {
		return
	}
	labelIndex, _ := GetLabelIndex(pod.Metadata.Labels)
	log.Info("GetLabelIndex: %v", labelIndex)
	// Step 1: Deal with new label
	if labelIndex == nil || len(labelIndex.Labels) == 0 {
		// a new label
		// create a new label index
		log.Info("create a new label index, %v", pod.Metadata.Labels)
		labelIndex = &api.LabelIndex{
			Labels:  pod.Metadata.Labels,
			PodName: []string{util.GetUniqueName(pod.Metadata.NameSpace, pod.Metadata.Name)},
		}

		// no service need to be updated, since the label is new
	} else {
		// update the label index
		labelIndex.PodName = append(labelIndex.PodName, util.GetUniqueName(pod.Metadata.NameSpace, pod.Metadata.Name))

		// need to update service
		for _, serviceName := range labelIndex.ServiceName {
			svc_namespace, name := util.GetNamespaceAndName(serviceName)
			svc, _ := GetService(svc_namespace, name)
			if svc != nil {
				// update service
				// add the new endpoint to the service

				// traverse all the service ports need to be exposed
				for _, targetPort := range svc.Spec.Ports {
					svc.Status.EndPoints = append(svc.Status.EndPoints,
						api.EndPoint{
							ServicePort: strconv.Itoa(targetPort.Port),
							IP:          pod.Status.PodIP,
							Ports:       matchTargetPort(targetPort, pod.Spec.Containers),
						})
				}
				log.Info("update service: %v", svc)
				err := UpdateService(svc)
				if err != nil {
					return
				}
			}
		}
	}
	// Step 2: store the new label index
	err := UpdateLabelIndex(labelIndex)
	if err != nil {
		log.Fatal("add label index error")
		return
	}
}

func OnPodUpdate(pod *api.Pod, oldLabel map[string]string) {
	log.Info("OnPodUpdate")
	if util.IsLabelEqual(pod.Metadata.Labels, oldLabel) {
		// no need to update
		return
	}
	labelIndex, _ := GetLabelIndex(pod.Metadata.Labels)
	log.Info("GetLabelIndex: %v", labelIndex)
	// Step 1: Deal with new label
	if labelIndex == nil || len(labelIndex.Labels) == 0 {
		// a new label
		// create a new label index
		log.Info("create a new label index, %v", pod.Metadata.Labels)
		labelIndex = &api.LabelIndex{
			Labels:  pod.Metadata.Labels,
			PodName: []string{util.GetUniqueName(pod.Metadata.NameSpace, pod.Metadata.Name)},
		}

		// no service need to be updated, since the label is new
	} else {
		// update the label index
		labelIndex.PodName = append(labelIndex.PodName, util.GetUniqueName(pod.Metadata.NameSpace, pod.Metadata.Name))

		// need to update service
		for _, serviceName := range labelIndex.ServiceName {
			svc_namespace, name := util.GetNamespaceAndName(serviceName)
			svc, _ := GetService(svc_namespace, name)
			if svc != nil {
				// update service
				// add the new endpoint to the service
				// traverse all the service ports need to be exposed
				for _, targetPort := range svc.Spec.Ports {
					svc.Status.EndPoints = append(svc.Status.EndPoints,
						api.EndPoint{
							ServicePort: strconv.Itoa(targetPort.Port),
							IP:          pod.Status.PodIP,
							Ports:       matchTargetPort(targetPort, pod.Spec.Containers),
						})
				}
				log.Info("update service: %v", svc)
				err := UpdateService(svc)
				if err != nil {
					return
				}
			}
		}
	}
	// Step 2: store the new label index
	err := UpdateLabelIndex(labelIndex)
	if err != nil {
		log.Fatal("add label index error")
		return
	}

	// Step 3: Deal with old label
	if oldLabel == nil {
		return
	}
	oldLabelIndex, _ := GetLabelIndex(oldLabel)
	if oldLabelIndex == nil || len(oldLabelIndex.Labels) == 0 {
		// Can't be here
		return
	}
	// remove the pod name from the old label index
	for i, name := range oldLabelIndex.PodName {
		if name == util.GetUniqueName(pod.Metadata.NameSpace, pod.Metadata.Name) {
			oldLabelIndex.PodName = append(oldLabelIndex.PodName[:i], oldLabelIndex.PodName[i+1:]...)
			break
		}
	}
	// update service
	for _, serviceName := range oldLabelIndex.ServiceName {
		svc_namespace, name := util.GetNamespaceAndName(serviceName)
		svc, _ := GetService(svc_namespace, name)
		if svc != nil {
			// update service
			// remove the old endpoint from the service
			for i, ep := range svc.Status.EndPoints {
				if ep.IP == pod.Status.PodIP {
					svc.Status.EndPoints = append(svc.Status.EndPoints[:i], svc.Status.EndPoints[i+1:]...)
					break
				}
			}
			err := UpdateService(svc)
			if err != nil {
				return
			}
		}
	}
	// store the old label index
	// check if to delete the label index
	if len(labelIndex.PodName) == 0 && len(labelIndex.ServiceName) == 0 {
		err := DeleteLabelIndex(pod.Metadata.Labels)
		if err != nil {
			log.Fatal("delete label index error")
		}
		return
	}
	err = UpdateLabelIndex(oldLabelIndex)
	if err != nil {
		log.Fatal("add label index error")
	}
}

func OnPodDelete(pod *api.Pod) {
	log.Info("OnPodDelete")
	labelIndex, _ := GetLabelIndex(pod.Metadata.Labels)
	if labelIndex == nil || len(labelIndex.Labels) == 0 {
		// Can't be here
		log.Error("Can't find labelIndex. %v", labelIndex)
		return
	}
	// remove the pod name from the label index
	for i, name := range labelIndex.PodName {
		if name == util.GetUniqueName(pod.Metadata.NameSpace, pod.Metadata.Name) {
			labelIndex.PodName = append(labelIndex.PodName[:i], labelIndex.PodName[i+1:]...)
			break
		}
	}
	// update service
	for _, serviceName := range labelIndex.ServiceName {
		svc_namespace, name := util.GetNamespaceAndName(serviceName)
		svc, _ := GetService(svc_namespace, name)
		if svc != nil {
			// update service
			// remove the old endpoint from the service
			for i, ep := range svc.Status.EndPoints {
				if ep.IP == pod.Status.PodIP {
					svc.Status.EndPoints = append(svc.Status.EndPoints[:i], svc.Status.EndPoints[i+1:]...)
					break
				}
			}
			err := UpdateService(svc)
			if err != nil {
				return
			}
		}
	}
	// store the label index
	// check if to delete the label index
	if len(labelIndex.PodName) == 0 && len(labelIndex.ServiceName) == 0 {
		err := DeleteLabelIndex(pod.Metadata.Labels)
		if err != nil {
			log.Fatal("delete label index error")
		}
		return
	}
	err := UpdateLabelIndex(labelIndex)
	if err != nil {
		log.Fatal("add label index error")
	}
}

func OnServiceUpdate(svc *api.Service, oldLabel map[string]string) {
	if util.IsLabelEqual(svc.Spec.Selector, oldLabel) {
		// no need to update
		return
	}

	// Step 1: Deal with new label
	labelIndex, _ := GetLabelIndex(svc.Spec.Selector)
	if labelIndex == nil || len(labelIndex.Labels) == 0 {
		// a new label
		// create a new label index
		labelIndex = &api.LabelIndex{
			Labels:      svc.Spec.Selector,
			ServiceName: []string{util.GetUniqueName(svc.Metadata.NameSpace, svc.Metadata.Name)},
		}
	} else {
		// update the label index
		labelIndex.ServiceName = append(labelIndex.ServiceName, util.GetUniqueName(svc.Metadata.NameSpace, svc.Metadata.Name))
		// update service's endpoint
	}
	for _, podName := range labelIndex.PodName {
		namespace, name := util.GetNamespaceAndName(podName)
		pod, _ := GetPod(namespace, name)
		if pod != nil {
			// traverse all the service ports need to be exposed
			for _, targetPort := range svc.Spec.Ports {
				svc.Status.EndPoints = append(svc.Status.EndPoints,
					api.EndPoint{
						ServicePort: strconv.Itoa(targetPort.Port),
						IP:          pod.Status.PodIP,
						Ports:       matchTargetPort(targetPort, pod.Spec.Containers),
					})
				log.Info("endpoint: %v", svc.Status.EndPoints)
			}
		}
	}
	// store service
	err := UpdateService(svc)
	if err != nil {
		return
	}
	// store the label index
	err = UpdateLabelIndex(labelIndex)
	if err != nil {
		log.Fatal("add label index error")
	}

	// Step 2: Deal with old label
	if oldLabel == nil {
		return
	}
	oldLabelIndex, _ := GetLabelIndex(oldLabel)
	if oldLabelIndex == nil || len(oldLabelIndex.Labels) == 0 {
		// Can't be here
		return
	}
	// remove the service name from the old label index
	for i, name := range oldLabelIndex.ServiceName {
		if name == util.GetUniqueName(svc.Metadata.NameSpace, svc.Metadata.Name) {
			oldLabelIndex.ServiceName = append(oldLabelIndex.ServiceName[:i], oldLabelIndex.ServiceName[i+1:]...)
			break
		}
	}
	// store the old label index
	// check if to delete the label index
	if len(oldLabelIndex.PodName) == 0 && len(oldLabelIndex.ServiceName) == 0 {
		err := DeleteLabelIndex(oldLabel)
		if err != nil {
			log.Fatal("delete label index error")
		}
		return
	}
	err = UpdateLabelIndex(oldLabelIndex)
	if err != nil {
		log.Fatal("add label index error")
	}

}

func OnServiceDelete(svc *api.Service) {
	labelIndex, _ := GetLabelIndex(svc.Spec.Selector)
	if labelIndex == nil || len(labelIndex.Labels) == 0 {
		// Can't be here
		return
	}
	// remove the service name from the label index
	for i, name := range labelIndex.ServiceName {
		if name == util.GetUniqueName(svc.Metadata.NameSpace, svc.Metadata.Name) {
			labelIndex.ServiceName = append(labelIndex.ServiceName[:i], labelIndex.ServiceName[i+1:]...)
			break
		}
	}
	// store the label index
	// check if to delete the label index
	if len(labelIndex.PodName) == 0 && len(labelIndex.ServiceName) == 0 {
		err := DeleteLabelIndex(svc.Spec.Selector)
		if err != nil {
			log.Fatal("delete label index error")
		}
		return
	}
	err := UpdateLabelIndex(labelIndex)
	if err != nil {
		log.Fatal("add label index error")
	}
}

func (e *EndPointController) Run() {
	log.Info("EndPointController is running")
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	topics := []string{msg_type.PodTopic, msg_type.ServiceTopic}
	e.subscriber.Subscribe(wg, ctx, topics, e)
	<-e.ready
	<-e.done
	cancel()
	wg.Wait()
}

// These functions interact with the API server

func GetService(namespace string, name string) (*api.Service, error) {
	URL := config.GetUrlPrefix() + config.ServiceURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, namespace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, name, -1)
	service := &api.Service{}
	err := httputil.Get(URL, service, "data")
	if err != nil {
		log.Error("err get service %s:%s", namespace, name)
		return nil, err
	}

	return service, nil
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

	err = httputil.Put(URL, serviceByteArray)

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

func matchTargetPort(targetPort api.ServicePort, containers []api.Container) []api.ContainerPort {
	matchedPorts := []api.ContainerPort{}

	// if targetPort not set, use the port
	if targetPort.TargetPort == "" {
		log.Info("targetPort not set")
		for _, container := range containers {
			for _, port := range container.Ports {
				if uint16(targetPort.Port) == uint16(port.ContainerPort) {
					matchedPorts = append(matchedPorts, port)
				}
			}
		}
		return matchedPorts
	}
	for _, container := range containers {
		for _, port := range container.Ports {
			// targetPort can be the name of the port or the number of the port
			log.Info("targetPort: %v, containerport: %v", targetPort.TargetPort, port.ContainerPort)
			if targetPort.TargetPort == port.Name || targetPort.TargetPort == strconv.Itoa(int(port.ContainerPort)) {
				matchedPorts = append(matchedPorts, port)
				log.Info("match target port: %v", port)
			}
		}
	}
	return matchedPorts
}
