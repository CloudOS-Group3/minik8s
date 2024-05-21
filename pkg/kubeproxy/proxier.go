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

// KubeProxy watch service changes
//
//	and update the ipvs rules
type KubeProxy struct {
	subscriber *kafka.Subscriber
	ready      chan bool
	done       chan bool
}

func (e KubeProxy) Setup(session sarama.ConsumerGroupSession) error {
	close(e.ready)
	return nil
}

func (e KubeProxy) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (e KubeProxy) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Info("Watch msg: %s\n", string(msg.Value))
		if msg.Topic == msg_type.ServiceTopic {
			session.MarkMessage(msg, "")
			serviceMsg := &msg_type.ServiceMsg{}
			err := json.Unmarshal(msg.Value, serviceMsg)
			if err != nil {
				log.Error("unmarshal service error")
				continue
			}
			switch serviceMsg.Opt {
			case msg_type.Update:
				// discard service without endpoints
				if len(serviceMsg.NewService.Status.EndPoints) == 0 {
					return nil
				}
				err := ipvs.UpdateService(&serviceMsg.NewService)
				if err != nil {
					log.Fatal("Failed to update service: %s", err.Error())
					return err
				}
				break
			case msg_type.Delete:
				err := ipvs.DeleteService(&serviceMsg.OldService)
				if err != nil {
					log.Fatal("Failed to delete service: %s", err.Error())
					return err
				}
				break
			case msg_type.Add:
				// discard service without endpoints
				if len(serviceMsg.NewService.Status.EndPoints) == 0 {
					return nil
				}
				err := ipvs.AddService(&serviceMsg.NewService)
				if err != nil {
					log.Fatal("Failed to add service: %s", err.Error())
					return err
				}
				break
			}
		}
	}
	return nil
}

func NewKubeProxy() *KubeProxy {
	brokers := []string{"127.0.0.1:9092"}
	group := "kube-proxy"
	return &KubeProxy{
		ready:      make(chan bool),
		done:       make(chan bool),
		subscriber: kafka.NewSubscriber(brokers, group),
	}
}

func (k *KubeProxy) Run() {
	log.Info("KubeProxy start")
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	topics := []string{msg_type.ServiceTopic}
	k.subscriber.Subscribe(wg, ctx, topics, k)
	<-k.ready
	<-k.done
	cancel()
	wg.Wait()
}
