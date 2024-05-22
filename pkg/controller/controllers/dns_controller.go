package controllers

import (
	"github.com/IBM/sarama"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/kafka"
)

type DNSController struct {
	RegisteredDNS []api.DNS
	ready         chan bool
	done          chan bool
	subscriber    *kafka.Subscriber
}

func NewDnsController() *DNSController {
	brokers := []string{"127.0.0.1:9092"}
	group := "dns-controller"
	Controller := &DNSController{
		ready:      make(chan bool),
		done:       make(chan bool),
		subscriber: kafka.NewSubscriber(brokers, group),
	}
	Controller.RegisteredDNS = make([]api.DNS, 0)
	return Controller
}

func (s *DNSController) Setup(_ sarama.ConsumerGroupSession) error {
	close(s.ready)
	return nil
}

func (s *DNSController) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (s *DNSController) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if msg.Topic == msg_type.DNSTopic {
			sess.MarkMessage(msg, "")
			s.DNSHandler(msg.Value)
		}
	}
	return nil
}

func (s *DNSController) DNSHandler(msg []byte) {

}
