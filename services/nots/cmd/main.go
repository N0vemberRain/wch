package main

import (
	"context"
	"log"
	"os"
	"strings"
	"wch/services/nots/internal/consumer"
)

func main() {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(brokers) == 0 {
		log.Fatal("event publisher settings: KAFKA_BROKERS not found")
	}

	consumer := consumer.NewConsumer(brokers, "message_sent", "notification-group")
	ctx := context.Background()

	go consumer.Start(ctx)

	select {
	case <-ctx.Done():
		log.Println("Shutting down...")
	}
}
