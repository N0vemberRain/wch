package main

import (
	"flag"
	"log"

	"net"

	handler "wch/services/chats/internal/adapters/grpc"
	pg "wch/services/chats/internal/adapters/postgres"
	"wch/services/chats/internal/controller"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pbchats "wch/gen/chats/v1"
)

const serviceName = "users"

func main() {
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()

	log.Printf("Starting the %s service on port %d", serviceName, port)

	db, err := pg.NewPostgresDBFromEnv()
	if err != nil {
		log.Fatalf("Database settings: %s", err)
	}

	userRepo := pg.NewChatRepository(db)
	userCtrl := controller.NewChatController(userRepo)
	h := handler.NewHandler(userCtrl)
	lis, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pbchats.RegisterChatsServiceServer(srv, h)
	reflection.Register(srv)
	srv.Serve(lis)
}
