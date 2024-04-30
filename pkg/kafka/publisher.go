package kafka

import (
	"github.com/IBM/sarama"
	"minik8s/util/log"
)

type Publisher struct {
	producer sarama.SyncProducer
}

func NewPublisher(addr []string) *Publisher {
	// make a default config (maybe add config as a parameter in the future)
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V3_6_0_0

	// start a producer
	producer, err := sarama.NewSyncProducer(addr, config)
	if err != nil {
		log.Fatal("Failed to start Sarama producer: %s", err.Error())
		return nil
	}
	return &Publisher{producer: producer}
}

func (p *Publisher) SendMessage(topic string, value string) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	_, _, err := p.producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}
