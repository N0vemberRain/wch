package main

import (
	"context"
	"log"
	"net"
	"os"
	"strings"
	chatspb "wch/gen/chats/v1"
	notspb "wch/gen/nots/v1"
	"wch/services/nots/internal/adapters"
	chatsgrpc "wch/services/nots/internal/adapters/chats_grpc"
	"wch/services/nots/internal/consumer"
	"wch/services/nots/internal/handler"
	"wch/services/nots/internal/service"
	"wch/services/nots/internal/streams"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const serviceName = "nots"

func main() {
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(brokers) == 0 {
		log.Fatal("event publisher settings: KAFKA_BROKERS not found")
	}

	chatsAddr, ok := os.LookupEnv("CHATS_SERVICE_ADDR")
	if !ok {
		log.Fatalf("Chats Service error: CHATS_SERVICE_ADDR is undefined")
	}

	chatsConn, err := grpc.NewClient(
		chatsAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Chats Service error: %s", err)
	}
	defer chatsConn.Close()

	chatsClient := chatspb.NewChatsServiceClient(chatsConn)

	chatProvider := chatsgrpc.NewChatProvider(chatsClient)

	hub := streams.NewStreamHub()
	sender := adapters.NewGRPCSender(hub)

	svc := service.NewNotificationService(sender, chatProvider)
	msgsHandler := adapters.NewMsgsHandlerGRPC(svc)

	consumer := consumer.NewConsumer(
		brokers,
		"message_sent",
		"notification-group",
		msgsHandler,
	)
	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	ctx := context.Background()

	go consumer.Start(ctx)

	// start gRPC server
	lis, _ := net.Listen("tcp", ":8080")
	srv := grpc.NewServer()

	grpcServer := handler.NewGRPCServer(hub)
	notspb.RegisterNotificationServiceServer(
		srv,
		grpcServer,
	)

	url, ok := os.LookupEnv("SERVICE_URL")
	if !ok {
		log.Fatalf("failed to listen: SERVICE_URL is undefined")
	}
	log.Printf("Starting the %s service: %s", serviceName, url)
	srv.Serve(lis)
}
