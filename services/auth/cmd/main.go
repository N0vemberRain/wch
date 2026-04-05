package main

import (
	"fmt"
	"log"
	"os"

	"net"

	"wch/pkg/auth"
	"wch/services/auth/internal/adapters/bcrypt"
	handler "wch/services/auth/internal/adapters/grpc"
	"wch/services/auth/internal/adapters/jwt"
	pg "wch/services/auth/internal/adapters/postgres"
	"wch/services/auth/internal/controller"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authpb "wch/gen/auth/v1"
)

const serviceName = "auth"

func main() {
	// var port int
	// flag.IntVar(&port, "port", 8087, "API handler port")
	// flag.Parse()

	// log.Printf("Starting the %s service on port %d", serviceName, port)

	db, err := pg.NewPostgresDBFromEnv()
	if err != nil {
		log.Fatalf("Database settings: %s", err)
	}

	authRepo := pg.NewCredentialsRepository(db)
	hasher := bcrypt.NewPasswordHasher(1)

	config, err := auth.LoadConfig()
	if err != nil {
		log.Fatalf("load auth config: %s\n", err.Error())
	}
	fmt.Printf("DURATION TIME: %v\n", config.ExpireTime)
	tokenValidator := auth.NewTokenValidator(config)
	tokenIssuer := jwt.NewTokenIssuer(config)

	authCtrl := controller.NewAuthController(authRepo, hasher, tokenIssuer)
	authHandler := handler.NewHandler(authCtrl)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthInterceptor(tokenValidator)),
	)
	authpb.RegisterAuthServiceServer(srv, authHandler)
	reflection.Register(srv)
	url, ok := os.LookupEnv("SERVICE_URL")
	if !ok {
		log.Fatalf("failed to listen: SERVICE_URL is undefined")
	}
	log.Printf("Starting the %s service: %s", serviceName, url)

	srv.Serve(lis)
}
