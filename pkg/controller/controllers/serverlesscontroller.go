package controllers

import (
	"context"
	"github.com/IBM/sarama"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/kafka"
	"sync"
)

type ServerlessController struct {
	jobs       []api.Job
	pods       []api.Pod
	subscriber *kafka.Subscriber
	ready      chan bool
	done       chan bool
}

func (this *ServerlessController) Setup(_ sarama.ConsumerGroupSession) error {
	close(this.ready)
	return nil
}

func (this *ServerlessController) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (this *ServerlessController) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if msg.Topic == msg_type.TriggerTopic {
			sess.MarkMessage(msg, "")
			this.triggerNewJob(msg.Value)
		} else if msg.Topic == msg_type.JobTopic {
			sess.MarkMessage(msg, "")
			this.updateJob(msg.Value)
		}
	}
	return nil
}

func (this *ServerlessController) triggerNewJob(content []byte) {

}

func (this *ServerlessController) updateJob(content []byte) {

}

func (this *ServerlessController) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	topics := []string{msg_type.TriggerTopic, msg_type.JobTopic}
	this.subscriber.Subscribe(wg, ctx, topics, this)
	<-this.ready
	<-this.done
	cancel()
	wg.Wait()
}
