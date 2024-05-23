package controllers

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"minik8s/pkg/api"
	"minik8s/pkg/api/msg_type"
	"minik8s/pkg/config"
	"minik8s/pkg/kafka"
	"minik8s/util/httputil"
	"minik8s/util/log"
	"strings"
	"sync"
	"time"
)

type NodeController struct {
	RegisteredNode []api.Node
	ready          chan bool
	done           chan bool
	subscriber     *kafka.Subscriber
}

func NewNodeController() *NodeController {
	brokers := []string{"127.0.0.1:9092"}
	group := "node-controller"
	Controller := &NodeController{
		ready:      make(chan bool),
		done:       make(chan bool),
		subscriber: kafka.NewSubscriber(brokers, group),
	}
	Controller.RegisteredNode = make([]api.Node, 0)
	return Controller
}

func (s *NodeController) Setup(_ sarama.ConsumerGroupSession) error {
	close(s.ready)
	return nil
}

func (s *NodeController) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (s *NodeController) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if msg.Topic == msg_type.NodeTopic {
			sess.MarkMessage(msg, "")
			s.NodeHandler(msg.Value)
		}
	}
	return nil
}

func (s *NodeController) CheckNode() {
	for {
		for index, node := range s.RegisteredNode {
			if node.Status.Condition.Status != api.NodeReady {
				continue
			}
			heartBeatTime := node.Status.Condition.LastHeartbeatTime
			currentTime := time.Now()
			TimeDiff := currentTime.Sub(heartBeatTime)
			if TimeDiff > time.Minute*3 {
				log.Info("Node dead: %s", node.Metadata.Name)
				node.Status.Condition.Status = api.NodeUnknown
				URL := config.GetUrlPrefix() + config.NodeURL
				URL = strings.Replace(URL, config.NamePlaceholder, node.Metadata.Name, -1)
				byteArr, err := json.Marshal(node)
				if err != nil {
					panic(err)
				}
				err = httputil.Put(URL, byteArr)
				if err != nil {
					panic(err)
				}
				s.RegisteredNode[index] = node
			}
		}
		time.Sleep(time.Second * 30)
	}
}

func (s *NodeController) NodeHandler(msg []byte) {
	var message msg_type.NodeMsg
	var node api.Node
	err := json.Unmarshal(msg, &message)
	if err != nil {
		panic(err)
	}
	if message.Opt == msg_type.Delete {
		for index, nodeInList := range s.RegisteredNode {
			if nodeInList.Metadata.Name == message.OldNode.Metadata.Name {
				s.RegisteredNode = append(s.RegisteredNode[:index], s.RegisteredNode[index+1:]...)
				return
			}
		}
	}
	node = message.NewNode
	exist := false
	for index, nodeInList := range s.RegisteredNode {
		if nodeInList.Metadata.Name == node.Metadata.Name {
			s.RegisteredNode[index] = node
			exist = true
		}
	}
	if !exist {
		log.Info("add node: %s", node.Metadata.Name)
		s.RegisteredNode = append(s.RegisteredNode, node)
	}
}

func (s *NodeController) Run() {
	go s.CheckNode()
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	topics := []string{msg_type.NodeTopic}
	s.subscriber.Subscribe(wg, ctx, topics, s)
	<-s.ready
	<-s.done
	cancel()
	wg.Wait()
}
