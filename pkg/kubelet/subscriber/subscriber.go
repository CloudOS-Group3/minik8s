package subscriber

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/kafka"
	"minik8s/pkg/kubelet/pod"
	"minik8s/util/log"
	"sync"
)

type KubeletSubscriber struct {
	subscriber *kafka.Subscriber
	pm         *pod.PodManager
	ready      chan bool
	done       chan bool
}

func NewKubeletSubscriber() *KubeletSubscriber {
	brokers := []string{"127.0.0.1:9092"}
	group := "kubelet"
	return &KubeletSubscriber{
		ready:      make(chan bool),
		done:       make(chan bool),
		pm:         pod.NewPodManager(),
		subscriber: kafka.NewSubscriber(brokers, group),
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
		if msg.Topic == "pod" {
			sess.MarkMessage(msg, "")
			k.PodHandler(msg.Value)
		}
	}
	return nil
}

func (k *KubeletSubscriber) PodHandler(msg []byte) {
	var podMsg msg_type.PodMsg
	err := json.Unmarshal(msg, &podMsg)
	if err != nil {
		log.Error("unmarshal pod message failed, error: %s", err.Error())
		panic(err)
	}
	switch podMsg.Opt {
	case msg_type.Add:
		k.pm.CreatePod(&podMsg.NewPod)
		break
	case msg_type.Delete:
		k.pm.DeletePod(&podMsg.OldPod)
		break
	case msg_type.Update:
		k.pm.DeletePod(&podMsg.OldPod)
		k.pm.CreatePod(&podMsg.NewPod)
		break
	}
}

func (k *KubeletSubscriber) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	topics := []string{"pod"}
	k.subscriber.Subscribe(wg, ctx, topics, k)
	<-k.ready
	<-k.done
	cancel()
	wg.Wait()
}
