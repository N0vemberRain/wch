package main

import (
	"flag"
	"log"

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
	var port int
	flag.IntVar(&port, "port", 8087, "API handler port")
	flag.Parse()

	log.Printf("Starting the %s service on port %d", serviceName, port)

	db, err := pg.NewPostgresDBFromEnv()
	if err != nil {
		log.Fatalf("Database settings: %s", err)
	}

	authRepo := pg.NewCredentialsRepository(db)
	hasher := bcrypt.NewPasswordHasher(1)
	issuer := jwt.NewTokenIssuer("secret", 15000)
	validator := auth.NewTokenValidator("secret", 15000)
	authCtrl := controller.NewAuthController(authRepo, hasher, issuer)
	authHandler := handler.NewHandler(authCtrl)

	lis, err := net.Listen("tcp", "localhost:8087")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthInterceptor(validator)),
	)
	authpb.RegisterAuthServiceServer(srv, authHandler)
	reflection.Register(srv)
	srv.Serve(lis)
}
