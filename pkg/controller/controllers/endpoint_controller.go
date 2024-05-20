package controllers

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"minik8s/pkg/api"
	msg_type "minik8s/pkg/api/msg_type"
	"minik8s/pkg/kafka"
	"minik8s/pkg/util"
	"minik8s/util/log"
	"sync"
)

type EndPointController struct {
	subscriber *kafka.Subscriber
	ready      chan bool
	done       chan bool
}

func (e EndPointController) Setup(session sarama.ConsumerGroupSession) error {
	close(e.ready)
	return nil
}

func (e EndPointController) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (e EndPointController) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Info("Watch msg: %s\n", string(msg.Value))
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
				if !util.IsLabelEqual(podMsg.NewPod.Spec.NodeSelector, podMsg.OldPod.Spec.NodeSelector) {
					OnPodUpdate(&podMsg.NewPod, podMsg.OldPod.Spec.NodeSelector)
				}
				break
			case msg_type.Delete:
				OnPodDelete(&podMsg.NewPod)
				break
			case msg_type.Add:
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
			switch serviceMsg.Opt {
			case msg_type.Update:
				if !util.IsLabelEqual(serviceMsg.NewService.Metadata.Labels, serviceMsg.OldService.Metadata.Labels) {
					OnServiceUpdate(&serviceMsg.NewService, serviceMsg.OldService.Metadata.Labels)
				}
				break
			case msg_type.Delete:
				OnServiceDelete(&serviceMsg.NewService)
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
	brokers := []string{"127.0.0.1:9092"}
	group := "endpoint-controller"
	return &EndPointController{
		ready:      make(chan bool),
		done:       make(chan bool),
		subscriber: kafka.NewSubscriber(brokers, group),
	}
}

func OnPodUpdate(pod *api.Pod, oldLabel map[string]string) {
	if util.IsLabelEqual(pod.Spec.NodeSelector, oldLabel) {
		// no need to update
		return
	}
	labelIndex, _ := GetLabelIndex(pod.Spec.NodeSelector)

	// Step 1: Deal with new label
	if labelIndex == nil {
		// a new label
		// create a new label index
		labelIndex = &api.LabelIndex{
			Labels:  pod.Spec.NodeSelector,
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
				var all_ports []api.ContainerPort
				for _, container := range pod.Spec.Containers {
					all_ports = append(all_ports, container.Ports...)
				}
				svc.Status.EndPoints = append(svc.Status.EndPoints,
					api.EndPoint{
						IP:    pod.Status.PodIP,
						Ports: all_ports,
					})
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
	if oldLabelIndex == nil {
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
		err := DeleteLabelIndex(pod.Spec.NodeSelector)
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
	labelIndex, _ := GetLabelIndex(pod.Spec.NodeSelector)
	if labelIndex == nil {
		// Can't be here
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
		err := DeleteLabelIndex(pod.Spec.NodeSelector)
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
	if util.IsLabelEqual(svc.Metadata.Labels, oldLabel) {
		// no need to update
		return
	}

	// Step 1: Deal with new label
	labelIndex, _ := GetLabelIndex(svc.Metadata.Labels)
	if labelIndex == nil {
		// Can't be here
		return
	}
	// update the label index
	labelIndex.ServiceName = append(labelIndex.ServiceName, util.GetUniqueName(svc.Metadata.NameSpace, svc.Metadata.Name))
	// update service's endpoint
	for _, podName := range labelIndex.PodName {
		namespace, name := util.GetNamespaceAndName(podName)
		pod, _ := GetPod(namespace, name)
		if pod != nil {
			var all_ports []api.ContainerPort
			for _, container := range pod.Spec.Containers {
				all_ports = append(all_ports, container.Ports...)
			}
			svc.Status.EndPoints = append(svc.Status.EndPoints,
				api.EndPoint{
					IP:    pod.Status.PodIP,
					Ports: all_ports,
				})
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
	if oldLabelIndex == nil {
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
	labelIndex, _ := GetLabelIndex(svc.Metadata.Labels)
	if labelIndex == nil {
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
		err := DeleteLabelIndex(svc.Metadata.Labels)
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
