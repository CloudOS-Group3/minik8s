package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/IBM/sarama"
	"sync"
	"testing"
	"time"
)

type TestConsumerGroupHandler struct{}

func (TestConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (TestConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h TestConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf("msg: %s\n", string(msg.Value))
		fmt.Printf("topic: %s\n", msg.Topic)
		fmt.Printf("partition: %d\n", msg.Partition)
		fmt.Printf("offset: %d\n", msg.Offset)
		if msg.Value == nil || len(msg.Value) == 0 || string(msg.Value) != "hello world" {
			return errors.New("failed to consume message")
		}
		sess.MarkMessage(msg, "")
		break
	}
	return nil
}

func TestKafka(t *testing.T) {
	brokers := []string{"127.0.0.1:9092"}
	consumer_group := "test-group"
	consumer_group_2 := "test-group-2"
	topics := []string{"test-topic"}
	producer, err := NewPublisher(brokers)
	if err != nil {
		t.Errorf("kafka producer init fail: %s", err)
	}
	defer producer.Close()
	consumer, err := NewSubscriber(brokers, consumer_group)
	consumer_2, err := NewSubscriber(brokers, consumer_group_2)
	if err != nil {
		t.Errorf("kafka consumer init fail: %s", err)
	}
	defer consumer_2.Close()
	defer consumer.Close()
	handler := TestConsumerGroupHandler{}
	fmt.Println("start consume")
	ctx, cancel := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			fmt.Printf("consuming\n")
			err := consumer.Consume(ctx, topics, handler)
			if err != nil {
				switch err {
				case sarama.ErrClosedClient, sarama.ErrClosedConsumerGroup:
					// 退出
					fmt.Printf("quit: kafka consumer")
					return
				case sarama.ErrOutOfBrokers:
					t.Errorf("kafka 崩溃了~")
				default:
					t.Errorf("kafka exception: %s", err.Error())
				}
				time.Sleep(1 * time.Second)
			} else {
				return
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			fmt.Printf("consuming\n")
			err := consumer_2.Consume(ctx2, topics, handler)
			if err != nil {
				switch err {
				case sarama.ErrClosedClient, sarama.ErrClosedConsumerGroup:
					// 退出
					fmt.Printf("quit: kafka consumer")
					return
				case sarama.ErrOutOfBrokers:
					t.Errorf("kafka 崩溃了~")
				default:
					t.Errorf("kafka exception: %s", err.Error())
				}
				time.Sleep(1 * time.Second)
			} else {
				return
			}
		}
	}()
	msg := &sarama.ProducerMessage{
		Topic: "test-topic",
		Value: sarama.ByteEncoder("hello world"),
	}
	fmt.Println("msg delivered")
	partition, offset, err := producer.SendMessage(msg)
	fmt.Printf("partition: %d, offset: %d\n", partition, offset)
	if err != nil {
		t.Errorf("kafka send fail: %s", err)
	}
	if err != nil {
		t.Errorf("kafka consumer fail: %s", err)
	}
	fmt.Println("start waiting")
	wg.Wait()
	cancel()
	cancel2()
	fmt.Println("stop consume")
}
