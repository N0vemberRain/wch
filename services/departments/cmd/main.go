package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	//"net/http"
	"net"
	"time"

	"wch/pkg/discovery"
	discmemory "wch/pkg/discovery/memory"
	"wch/services/departments/internal/controller"
	"wch/services/departments/internal/repository"

	//mdgateway "wch/users/internal/gateway/metadata/http"
	//httphandler "wch/users/internal/handler/http"

	"wch/gen"
	//"wch/pkg/discovery/memory"

	grpchandler "wch/services/departments/internal/handler/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "users"

func main() {
	var port int
	flag.IntVar(&port, "port", 8083, "API handler port")
	flag.Parse()

	log.Printf("Starting the rating service on port %d\n", port)
	//registry := discmemory.NewRegistry("localhost:8500")
	registry := discmemory.NewRegistry()
	//if err := nil {
	//panic(err)
	//}
	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(
		ctx, instanceID,
		serviceName,
		fmt.Sprintf("localhost:%d", port),
	); err != nil {
		panic(err)
	}

	go func() {
		for {
			err := registry.ReportHealthState(instanceID, serviceName)
			if err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()

	defer registry.Deregister(ctx, instanceID, serviceName)

	db, err := repository.NewPostgresDBFromEnv()
	if err != nil {
		log.Fatalf("Database settings: %s", err)
	}
	departRepo := repository.NewRepositoryPg(db)
	depsCtrl := controller.NewController(departRepo)

	h := grpchandler.New(depsCtrl)
	lis, err := net.Listen("tcp", "localhost:8083")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	gen.RegisterDepsServiceServer(srv, h)
	reflection.Register(srv)
	srv.Serve(lis)
}
