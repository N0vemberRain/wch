package main

import (
	"flag"
	"log"

	//"net/http"
	"net"

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
	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()

	log.Println("Starting the rating service on port %d", port)
	//registry := discmemory.NewRegistry("localhost:8500")
	//registry := discmemory.NewRegistry()
	//if err := nil {
	//panic(err)
	//}
	// ctx := context.Background()
	// instanceID := discovery.GenerateInstanceID(serviceName)

	// if err := registry.Register(
	// 	ctx, instanceID,
	// 	serviceName,
	// 	fmt.Sprintf("localhost:%d", port),
	// ); err != nil {
	// 	panic(err)
	// }

	// go func() {
	// 	for {
	// 		err := registry.ReportHealthState(instanceID, serviceName)
	// 		if err != nil {
	// 			log.Println("Failed to report healthy state: " + err.Error())
	// 		}
	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	// defer registry.Deregister(ctx, instanceID, serviceName)

	db, err := pg.NewPostgresDBFromEnv()
	if err != nil {
		log.Fatalf("Database settings: %s", err)
	}

	userRepo := pg.NewUserRepositoryPg(db)
	userCtrl := controller.NewUserController(userRepo)
	h := grpchandler.New(userCtrl)
	lis, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	userspb.RegisterUsersServiceServer(srv, h)
	reflection.Register(srv)
	srv.Serve(lis)
}
