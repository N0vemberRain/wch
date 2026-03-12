package main

import (
	"log"
	"os"

	//"net/http"
	"net"

	"wch/pkg/auth"
	"wch/services/users/internal/controller"

	//mdgateway "wch/users/internal/gateway/metadata/http"
	//httphandler "wch/users/internal/handler/http"

	userspb "wch/gen/users/v1"
	//"wch/pkg/discovery/memory"

	grpchandler "wch/services/users/internal/adapters/grpc"
	pg "wch/services/users/internal/adapters/postgres"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "users"

func main() {
	// var port int
	// flag.IntVar(&port, "port", 8082, "API handler port")
	// flag.Parse()

	// log.Printf("Starting the rating service on port %d\n", port)

	db, err := pg.NewPostgresDBFromEnv()
	if err != nil {
		log.Fatalf("Database settings: %s", err)
	}

	userRepo := pg.NewUserRepositoryPg(db)
	userCtrl := controller.NewUserController(userRepo)
	h := grpchandler.New(userCtrl)
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
	userspb.RegisterUsersServiceServer(srv, h)
	reflection.Register(srv)
	url, ok := os.LookupEnv("SERVICE_URL")
	if !ok {
		log.Fatalf("failed to listen: SERVICE_URL is undefined")
	}
	log.Printf("Starting the %s service: %s", serviceName, url)
	srv.Serve(lis)
}
