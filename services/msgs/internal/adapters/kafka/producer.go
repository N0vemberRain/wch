package kafka

import (
	"context"
	"fmt"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(brokers []string, topic string) (*Producer, error) {
	writer := kafkago.Writer{
		Addr:                   kafkago.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafkago.LeastBytes{},
		AllowAutoTopicCreation: true,
		Transport: &kafkago.Transport{
			MetadataTTL: 500 * time.Millisecond,
		},
		MaxAttempts: 5,
	}

	if writer.Addr == nil {
		var err_msg string
		for _, str := range brokers {

			err_msg = err_msg + str + " "
		}
		return nil, fmt.Errorf("broker's format is invalid: %s", err_msg)
	}

	return &Producer{
		writer: &writer,
	}, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, key string, value []byte) error {
	fmt.Printf("kafka.Producer.Publish: write messages: %s\n", topic)
	return p.writer.WriteMessages(ctx, kafkago.Message{
		Key:   []byte(key),
		Value: value,
	})
}
