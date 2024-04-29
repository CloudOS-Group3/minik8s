package kafka

import (
	"github.com/IBM/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

func NewSubscriber(brokers []string, group string, topics []string) (*cluster.Consumer, error) {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	config.Group.Return.Notifications = true
	return cluster.NewConsumer(brokers, group, topics, config)
}
