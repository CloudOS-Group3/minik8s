package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"minik8s/pkg/api"
	"minik8s/pkg/kafka"
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
	// TODO: require node list from apiserver
	nodeList := make([]api.Node, 2)
	nodeList[0].Metadata.Name = "node1"
	nodeList[1].Metadata.Name = "node2"

	brokers := []string{"127.0.0.1:9092"}
	group := "scheduler"
	return &Scheduler{
		nodes:      nodeList,
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
		if msg.Topic == "pod" {
			sess.MarkMessage(msg, "")
			s.PodHandler(msg.Value)
		} else if msg.Topic == "node" {
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
	// TODO: send new node to apiserver
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
	topics := []string{"pod", "node"}
	s.subscriber.Subscribe(wg, ctx, topics, s)
	<-s.ready
	<-s.done
	cancel()
	wg.Wait()
}
