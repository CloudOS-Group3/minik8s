package kafka

import (
	"github.com/IBM/sarama"
)

func NewSubscriber(brokers []string, group string) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Consumer.Return.Errors = true

	return sarama.NewConsumerGroup(brokers, group, config)
}
