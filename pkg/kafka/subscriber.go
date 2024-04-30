package kafka

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"minik8s/util/log"
	"sync"
	"time"
)

type Subscriber struct {
	consumerGroup sarama.ConsumerGroup
}

func NewSubscriber(brokers []string, group string) *Subscriber {
	// make a default config (maybe add config as a parameter in the future)
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	// create a consumer group
	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Fatal("Error creating consumer group: %s", err.Error())
		return nil
	}
	return &Subscriber{consumerGroup: consumerGroup}
}

func (s *Subscriber) Subscribe(wg *sync.WaitGroup, ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) {
	/* usage (not pretty sure):
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	Subscribe(...)
	<- handler.ready (handler.setup should close ready)
	wg.Wait()
	cancel()
	*/
	// use go func() to run it async
	// maybe we should give a way to terminate it
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := s.consumerGroup.Consume(ctx, topics, handler)
			if err != nil {
				switch err {
				case sarama.ErrClosedClient, sarama.ErrClosedConsumerGroup:
					// kafka consumer quit
					fmt.Printf("quit: kafka consumer\n")
					return
				case sarama.ErrOutOfBrokers:
					fmt.Printf("kafka crash\n")
				default:
					fmt.Printf("kafka exception: %s\n", err.Error())
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()
}
