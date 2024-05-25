package subscriber

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/kafka"
	"minik8s/pkg/kubelet/host"
	"minik8s/pkg/kubelet/node"
	"minik8s/pkg/kubelet/pod"
	"minik8s/util/log"
	"sync"
)

type KubeletSubscriber struct {
	subscriber  *kafka.Subscriber
	HostManager *host.KubeletHostManager
	ready       chan bool
	done        chan bool
}

func NewKubeletSubscriber() *KubeletSubscriber {
	group := "kubelet" + "-" + config.Nodename
	return &KubeletSubscriber{
		ready:       make(chan bool),
		done:        make(chan bool),
		subscriber:  kafka.NewSubscriber(group),
		HostManager: host.NewHostManager(),
	}
}

func (k *KubeletSubscriber) Setup(_ sarama.ConsumerGroupSession) error {
	close(k.ready)
	return nil
}

func (k *KubeletSubscriber) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (k *KubeletSubscriber) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Info("Message claimed: value %s", string(msg.Value))
		if msg.Topic == msg_type.PodTopic {
			sess.MarkMessage(msg, "")
			k.PodHandler(msg.Value)
		}
		if msg.Topic == msg_type.DNSTopic {
			sess.MarkMessage(msg, "")
			k.DNSHandler(msg.Value)
		}
	}
	return nil
}

func (k *KubeletSubscriber) PodHandler(msg []byte) {
	var podMsg msg_type.PodMsg
	err := json.Unmarshal(msg, &podMsg)
	if err != nil {
		log.Error("unmarshal pod message failed, error: %s", err.Error())
		return
	}
	switch podMsg.Opt {
	case msg_type.Add:
		if podMsg.NewPod.Spec.NodeName != config.Nodename {
			break
		}
		log.Info("create pod: %v", podMsg.NewPod)
		pod.CreatePod(&podMsg.NewPod)
		break
	case msg_type.Delete:
		if podMsg.OldPod.Spec.NodeName != config.Nodename {
			break
		}
		log.Info("delete pod: %v", podMsg.OldPod)
		pod.DeletePod(&podMsg.OldPod)
		break
	case msg_type.Update:
		OldSpec, _ := json.Marshal(podMsg.OldPod.Spec)
		NewSpec, _ := json.Marshal(podMsg.NewPod.Spec)
		if string(OldSpec) != string(NewSpec) {
			log.Info("update pod: %v", podMsg.NewPod)
			if podMsg.OldPod.Spec.NodeName == config.Nodename {
				pod.DeletePod(&podMsg.OldPod)
			}
			if podMsg.NewPod.Spec.NodeName == config.Nodename {
				pod.CreatePod(&podMsg.NewPod)
			}
		}
		break
	}
}

func (k *KubeletSubscriber) DNSHandler(msg []byte) {
	var dnsMsg msg_type.DNSMsg
	err := json.Unmarshal(msg, &dnsMsg)
	if err != nil {
		log.Error("unmarshal dns message failed, error: %s", err.Error())
		return
	}
	switch dnsMsg.Opt {
	case msg_type.Add:
		log.Info("Add DNS: %s", dnsMsg.NewDNS.Host)
		k.HostManager.AddHost(dnsMsg.NewDNS.Host)
		break
	case msg_type.Delete:
		log.Info("Delete DNS: %s", dnsMsg.OldDNS.Host)
		k.HostManager.RemoveHost(dnsMsg.OldDNS.Host)
		break
	case msg_type.Update:
		if dnsMsg.OldDNS.Host != dnsMsg.NewDNS.Host {
			k.HostManager.RemoveHost(dnsMsg.OldDNS.Host)
			k.HostManager.AddHost(dnsMsg.NewDNS.Host)
		}
	}
}

func (k *KubeletSubscriber) Run() {
	go node.DoHeartBeat()
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	topics := []string{msg_type.PodTopic, msg_type.DNSTopic}
	k.subscriber.Subscribe(wg, ctx, topics, k)
	<-k.ready
	<-k.done
	cancel()
	wg.Wait()
}
