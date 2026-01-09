package main

import (
	"flag"
	"log"

	"net"

	handler "wch/services/msgs/internal/adapters/grpc"
	pg "wch/services/msgs/internal/adapters/postgres"

	"wch/services/msgs/internal/controller"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	msgspb "wch/gen/msgs/v1"
)

const serviceName = "msgs"

func main() {
	var port int
	flag.IntVar(&port, "port", 8086, "API handler port")
	flag.Parse()

	log.Printf("Starting the %s service on port %d", serviceName, port)

	db, err := pg.NewPostgresDBFromEnv()
	if err != nil {
		log.Fatalf("Database settings: %s", err)
	}

	msgsRepo := pg.NewMessageRepository(db)
	msgsCtrl := controller.NewMessageController(msgsRepo)
	h := handler.NewHandler(msgsCtrl)
	lis, err := net.Listen("tcp", "localhost:8086")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	msgspb.RegisterMessagesServiceServer(srv, h)
	reflection.Register(srv)
	srv.Serve(lis)

}
