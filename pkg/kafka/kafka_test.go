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

type TestConsumerGroupHandler struct {
	ready chan bool
	done  chan bool
}

func (h TestConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	// mark consumer setup
	close(h.ready)
	return nil
}
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
		close(h.done)
		sess.MarkMessage(msg, "")
		break
	}
	return nil
}

func TestKafka(t *testing.T) {
	brokers := []string{"127.0.0.1:9092"}
	consumer_group := "test-group"
	consumer_group_2 := "test-group-2"
	topics := []string{"test-topic-r"}
	publisher := NewPublisher(brokers)
	if publisher == nil {
		t.Errorf("kafka producer init fail")
		return
	}
	defer publisher.producer.Close()
	consumerGroup := NewSubscriber(brokers, consumer_group)
	consumerGroup_another := NewSubscriber(brokers, consumer_group_2)
	if consumerGroup == nil || consumerGroup_another == nil {
		t.Errorf("kafka consumer init fail")
		return
	}
	defer consumerGroup.consumerGroup.Close()
	defer consumerGroup_another.consumerGroup.Close()
	handler := TestConsumerGroupHandler{
		ready: make(chan bool),
		done:  make(chan bool),
	}

	handler2 := TestConsumerGroupHandler{
		ready: make(chan bool),
		done:  make(chan bool),
	}
	fmt.Println("start consume")
	ctx, cancel := context.WithCancel(context.Background())
	ctx2, cancel2 := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-handler.done:
				return
			default:
				fmt.Printf("consuming\n")
				err := consumerGroup.consumerGroup.Consume(ctx, topics, handler)
				if err != nil {
					switch err {
					case sarama.ErrClosedClient, sarama.ErrClosedConsumerGroup:
						fmt.Printf("quit: kafka consumer")
						return
					case sarama.ErrOutOfBrokers:
						t.Errorf("kafka crash")
					default:
						t.Errorf("kafka exception: %s", err.Error())
					}
					time.Sleep(1 * time.Second)
				} else {
					return
				}
			}

		}
	}()
	// wait for consumer setup
	<-handler.ready

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-handler2.done:
				return
			default:
				fmt.Printf("consuming\n")
				err := consumerGroup_another.consumerGroup.Consume(ctx2, topics, handler2)
				if err != nil {
					switch err {
					case sarama.ErrClosedClient, sarama.ErrClosedConsumerGroup:
						fmt.Printf("quit: kafka consumer")
						return
					case sarama.ErrOutOfBrokers:
						t.Errorf("kafka crash")
					default:
						t.Errorf("kafka exception: %s", err.Error())
					}
					time.Sleep(1 * time.Second)
				} else {
					return
				}
			}
		}
	}()
	<-handler2.ready
	for {
		select {
		case <-handler.done:
			select {
			case <-handler2.done:
				cancel()
				cancel2()
				wg.Wait()
				return
			}
		default:
			msg := &sarama.ProducerMessage{
				Topic: "test-topic-r",
				Value: sarama.ByteEncoder("hello world"),
			}
			fmt.Println("msg delivered")
			partition, offset, err := publisher.producer.SendMessage(msg)
			fmt.Printf("partition: %d, offset: %d\n", partition, offset)
			if err != nil {
				t.Errorf("kafka send fail: %s", err)
			}
			if err != nil {
				t.Errorf("kafka consumer fail: %s", err)
			}
			time.Sleep(3 * time.Second)
		}
	}
}
