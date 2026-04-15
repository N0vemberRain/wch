package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type MessageSentEvent struct {
	Type      string
	ID        uuid.UUID
	ChatID    uuid.UUID
	SenderID  uuid.UUID
	Content   string
	Timestamp time.Time
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, topic string, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers: brokers,
				Topic:   topic,
				GroupID: groupID,
			},
		),
	}
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Println("error reading message: ", err)
			continue
		}

		var event MessageSentEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("error parsing message: ", err)
			continue
		}

		log.Println("Received message: ", event)
	}
}
