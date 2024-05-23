package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/kafka"
	"minik8s/util/httputil"
	"strings"
	"sync"
)

type Scheduler struct {
	nodes      []api.Node
	ready      chan bool
	done       chan bool
	subscriber *kafka.Subscriber
	count      int
}

func NewScheduler() *Scheduler {
	URL := config.GetUrlPrefix() + config.NodesURL
	var initialNode []api.Node
	_ = httputil.Get(URL, &initialNode, "data")
	brokers := []string{"127.0.0.1:9092"}
	group := "scheduler"
	return &Scheduler{
		nodes:      initialNode,
		ready:      make(chan bool),
		done:       make(chan bool),
		count:      0,
		subscriber: kafka.NewSubscriber(brokers, group),
	}
}

func (s *Scheduler) Setup(_ sarama.ConsumerGroupSession) error {
	close(s.ready)
	return nil
}

func (s *Scheduler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (s *Scheduler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if msg.Topic == msg_type.PodTopic {
			sess.MarkMessage(msg, "")
			s.PodHandler(msg.Value)
		} else if msg.Topic == msg_type.NodeTopic {
			sess.MarkMessage(msg, "")
			s.NodeHandler(msg.Value)
		}
	}
	return nil
}

func (s *Scheduler) PodHandler(msg []byte) {
	var pod api.Pod
	err := json.Unmarshal(msg, &pod)
	if err != nil {
		panic(err)
	}
	if pod.Spec.NodeName != "" {
		return
	} else {
		index := s.count % len(s.nodes)
		s.count = s.count + 1
		pod.Spec.NodeName = s.nodes[index].Metadata.Name
	}
	fmt.Printf("pod %s has assigned to node %s\n", pod.Metadata.Name, pod.Spec.NodeName)

	URL := config.GetUrlPrefix() + config.PodURL
	URL = strings.Replace(URL, config.NamespacePlaceholder, pod.Metadata.NameSpace, -1)
	URL = strings.Replace(URL, config.NamePlaceholder, pod.Metadata.Name, -1)
	byteArr, err := json.Marshal(pod)
	if err != nil {
		panic(err)
	}
	err = httputil.Put(URL, byteArr)
}

func (s *Scheduler) NodeHandler(msg []byte) {
	var node api.Node
	err := json.Unmarshal(msg, &node)
	if err != nil {
		panic(err)
	}
	exist := false
	for index, nodeInList := range s.nodes {
		if nodeInList.Metadata.Name == node.Metadata.Name {
			s.nodes[index] = node
			exist = true
		}
	}
	if !exist {
		s.nodes = append(s.nodes, node)
	}
}

func (s *Scheduler) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	topics := []string{msg_type.PodTopic, msg_type.NodeTopic}
	s.subscriber.Subscribe(wg, ctx, topics, s)
	<-s.ready
	<-s.done
	cancel()
	wg.Wait()
}
