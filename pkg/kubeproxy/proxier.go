package kubeproxy

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/kafka"
	"minik8s/pkg/kubeproxy/ipvs"
	"minik8s/util/log"
	"sync"
)

type ProxyInterface interface {
	OnServiceCreate(service *api.Service) error
	OnServiceUpdate(oldService, newService *api.Service) error
	OnServiceDelete(service *api.Service) error
}

type KubeproxySub struct {
	subscriber *kafka.Subscriber
	ready      chan bool
	done       chan bool
}

func NewKubeproxySub() *KubeproxySub {
	brokers := []string{"127.0.0.1:9092"}
	group := "kubeproxy"
	return &KubeproxySub{
		ready:      make(chan bool),
		done:       make(chan bool),
		subscriber: kafka.NewSubscriber(brokers, group),
	}
}

func (k *KubeproxySub) Setup(_ sarama.ConsumerGroupSession) error {
	close(k.ready)
	return nil
}

func (k *KubeproxySub) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (k *KubeproxySub) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Info("Message claimed: value %s", string(msg.Value))
		if msg.Topic == msg_type.EndpointTopic {
			sess.MarkMessage(msg, "")
			k.EndpointHandler(msg.Value)
		}
	}
	return nil
}

func (k *KubeproxySub) EndpointHandler(msg []byte) {
	var serviceMsg msg_type.ServiceMsg
	err := json.Unmarshal(msg, &serviceMsg)
	if err != nil {
		log.Error("unmarshal pod message failed, error: %s", err.Error())
		panic(err)
	}
	switch serviceMsg.Opt {
	case msg_type.Add:
		log.Info("create pod: %v", serviceMsg.NewService)
		err := ipvs.AddService(&serviceMsg.NewService)
		if err != nil {
			log.Error("add service failed, error: %s", err.Error())
			return
		}
		break
	case msg_type.Delete:
		log.Info("delete pod: %v", serviceMsg.OldService)
		err := ipvs.DeleteService(&serviceMsg.OldService)
		if err != nil {
			log.Error("delete service failed, error: %s", err.Error())
			return
		}
		break
	case msg_type.Update:
		log.Info("update pod: %v", serviceMsg.NewService)
		err := ipvs.UpdateService(&serviceMsg.NewService)
		if err != nil {
			log.Error("update service failed, error: %s", err.Error())
			return
		}
	}
}

func (k *KubeproxySub) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	topics := []string{msg_type.EndpointTopic}
	k.subscriber.Subscribe(wg, ctx, topics, k)
	<-k.ready
	<-k.done
	cancel()
	wg.Wait()
}
