package kafka

import (
	"fmt"
	"github.com/IBM/sarama"
	"testing"
)

func TestKafka(t *testing.T) {
	brokers := []string{"127.0.0.1:9092"}
	consumer_group := "test-group"
	topics := []string{"test-topic"}
	producer, err := NewPublisher(brokers)
	if err != nil {
		t.Errorf("kafka producer init fail: %s", err)
	}
	defer producer.Close()
	consumer, err := NewSubscriber(brokers, consumer_group, topics)
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
	select {
	case msg, ok := <-consumer.Messages():
		if ok {
			if string(msg.Value) != "hello world" {
				t.Errorf("kafka receive wrong message: %s", string(msg.Value))
			}
			consumer.MarkOffset(msg, "") // 上报offset
		} else {
			fmt.Println("监听服务失败")
		}
	}
}
