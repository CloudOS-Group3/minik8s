package kafka

import (
	"context"
	"errors"
	"github.com/IBM/sarama"
	"testing"
)

type TestConsumerGroupHandler struct{}

func (TestConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (TestConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h TestConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		if msg.Value == nil || len(msg.Value) == 0 || string(msg.Value) != "hello world" {
			return errors.New("failed to consume message")
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func TestKafka(t *testing.T) {
	brokers := []string{"127.0.0.1:9092"}
	consumer_group := "test-group"
	topics := []string{"test-topic"}
	producer, err := NewPublisher(brokers)
	if err != nil {
		t.Errorf("kafka producer init fail: %s", err)
	}
	defer producer.Close()
	consumer, err := NewSubscriber(brokers, consumer_group)
	if err != nil {
		t.Errorf("kafka consumer init fail: %s", err)
	}
	defer consumer.Close()
	msg := &sarama.ProducerMessage{
		Topic: topics[0],
		Value: sarama.ByteEncoder("hello world"),
	}
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		t.Errorf("kafka send fail: %s", err)
	}
	ctx := context.Background()
	handler := TestConsumerGroupHandler{}
	err = consumer.Consume(ctx, topics, handler)
	if err != nil {
		t.Errorf("kafka consumer fail: %s", err)
	}
}
