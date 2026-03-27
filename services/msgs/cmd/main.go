package main

import (
	"log"
	"os"
	"strings"

	"net"

	"wch/pkg/auth"
	chatsgrpc "wch/services/msgs/internal/adapters/chats_grpc"
	handler "wch/services/msgs/internal/adapters/grpc"
	"wch/services/msgs/internal/adapters/kafka"
	pg "wch/services/msgs/internal/adapters/postgres"
	"wch/services/msgs/internal/domain/shared"

	"wch/services/msgs/internal/controller"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	chatspb "wch/gen/chats/v1"
	msgspb "wch/gen/msgs/v1"
)

const serviceName = "msgs"

func main() {
	// var port int
	// flag.IntVar(&port, "port", 8086, "API handler port")
	// flag.Parse()

	// log.Printf("Starting the %s service on port %d", serviceName, port)
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

	db, err := pg.NewPostgresDBFromEnv()
	if err != nil {
		log.Fatalf("Database settings: %s", err)
	}

	msgsRepo := pg.NewMessageRepository(db)

	chatsClient := chatspb.NewChatsServiceClient(chatsConn)
	chatsProvider := chatsgrpc.NewChatProvider(chatsClient)
	permChecker := shared.NewPermissionChecker(chatsProvider)
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	if len(brokers) == 0 {
		log.Fatal("event publisher settings: KAFKA_BROKERS not found")
	}
	eventPublisher, err := kafka.NewProducer(brokers, "message_sent")
	if err != nil {
		log.Fatalf("event publisher settings: %s", err.Error())
	}

	msgsCtrl := controller.NewMessageController(msgsRepo, permChecker, eventPublisher)
	h := handler.NewHandler(msgsCtrl)
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	config, err := auth.LoadConfig()
	if err != nil {
		log.Fatalf("load auth config: %s\n", err.Error())
	}

	tokenValidator := auth.NewTokenValidator(config)

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthInterceptor(tokenValidator)),
	)
	msgspb.RegisterMessagesServiceServer(srv, h)
	reflection.Register(srv)
	url, ok := os.LookupEnv("SERVICE_URL")
	if !ok {
		log.Fatalf("failed to listen: SERVICE_URL is undefined")
	}
	log.Printf("Starting the %s service: %s", serviceName, url)
	srv.Serve(lis)

}
