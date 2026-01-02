package main

import (
	"flag"
	"log"

	"net"

	handler "wch/services/chats/internal/adapters/grpc"
	pg "wch/services/chats/internal/adapters/postgres"
	up "wch/services/chats/internal/adapters/users_grpc"

	"wch/services/chats/internal/controller"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	pbchats "wch/gen/chats/v1"
	userspb "wch/gen/users/v1"
)

const serviceName = "users"

func main() {
	var port int
	flag.IntVar(&port, "port", 8085, "API handler port")
	flag.Parse()

	log.Printf("Starting the %s service on port %d", serviceName, port)

	usersConn, err := grpc.NewClient(
		"localhost:8082",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Users Service error: %s", err)
	}

	db, err := pg.NewPostgresDBFromEnv()
	if err != nil {
		log.Fatalf("Database settings: %s", err)
	}

	chatsRepo := pg.NewChatRepository(db)
	usersClient := userspb.NewUsersServiceClient(usersConn)
	userProvider := up.NewUserProvider(usersClient)
	chatsCtrl := controller.NewChatController(chatsRepo, userProvider)
	h := handler.NewHandler(chatsCtrl)
	lis, err := net.Listen("tcp", "localhost:8085")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pbchats.RegisterChatsServiceServer(srv, h)
	reflection.Register(srv)
	srv.Serve(lis)
}
