package kafka

import "github.com/IBM/sarama"

func NewPublisher(addr []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Version = sarama.V3_6_0_0
	return sarama.NewSyncProducer(addr, config)
}
