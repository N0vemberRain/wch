package consumer

import (
	"context"
	"encoding/json"
	"log"

	"wch/pkg/events"
	"wch/services/nots/internal/handler"

	"github.com/segmentio/kafka-go"
)

// type MessageSentEvent struct {
// 	Type      string
// 	ID        uuid.UUID
// 	ChatID    uuid.UUID
// 	SenderID  uuid.UUID
// 	Content   string
// 	Timestamp time.Time
// }

type Consumer struct {
	reader  *kafka.Reader
	handler handler.MsgsHandler
}

func NewConsumer(brokers []string, topic string, groupID string, handler handler.MsgsHandler) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(
			kafka.ReaderConfig{
				Brokers: brokers,
				Topic:   topic,
				GroupID: groupID,
			},
		),
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	log.Println("Consumer (msgs): starting")
	for {
		log.Println("Consumer (msgs): waiting new messages")
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			log.Println("error reading message: ", err)
			continue
		}

		var event events.MessageSentEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Println("error parsing message: ", err)
			continue
		}

		c.handler.Handle(ctx, event)
		log.Println("Received message: ", event)
	}
}
