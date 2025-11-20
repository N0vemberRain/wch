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
	"wch/services/users/internal/controller"
	"wch/services/users/internal/repository"

	//mdgateway "wch/users/internal/gateway/metadata/http"
	//httphandler "wch/users/internal/handler/http"

	"wch/gen"
	//"wch/pkg/discovery/memory"

	grpchandler "wch/services/users/internal/handler/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serviceName = "users"

func main() {
	// log.Println("Starting the users service")
	// registry, err := static.NewRegistry(map[string][]string{
	// 	"metadata": {"localhost:8081"},
	// 	"users":    {"localhost:8082"},
	// })

	// ctx := context.Background()
	// if err := registry.Register(ctx, "users", "localhost:8082"); err != nil {
	// 	panic(err)
	// }
	// defer registry.Deregister(ctx, "users")

	var port int
	flag.IntVar(&port, "port", 8082, "API handler port")
	flag.Parse()

	log.Println("Starting the rating service on port %d", port)
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
	departRepo := repository.NewDepartmentRepositoryPg(db)
	depsCtrl := controller.NewDepartmentController(departRepo)

	userRepo := repository.NewUserRepositoryPg(db)
	userCtrl := controller.NewUserController(userRepo)
	h := grpchandler.New(userCtrl, depsCtrl)
	lis, err := net.Listen("tcp", "localhost:8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	gen.RegisterUsersServiceServer(srv, h)
	reflection.Register(srv)
	srv.Serve(lis)
}

// func main() {
// 	var port int
// 	flag.IntVar(&port, "port", 8082, "API handler port")
// 	flag.Parse()

// 	log.Println("Starting the rating service on port %d", port)
// 	//registry := discmemory.NewRegistry("localhost:8500")
// 	registry := discmemory.NewRegistry()
// 	//if err := nil {
// 	//panic(err)
// 	//}
// 	ctx := context.Background()
// 	instanceID := discovery.GenerateInstanceID(serviceName)

// 	if err := registry.Register(
// 		ctx, instanceID,
// 		serviceName,
// 		fmt.Sprintf("localhost:%d", port),
// 	); err != nil {
// 		panic(err)
// 	}

// 	go func() {
// 		for {
// 			err := registry.ReportHealthState(instanceID, serviceName)
// 			if err != nil {
// 				log.Println("Failed to report healthy state: " + err.Error())
// 			}
// 			time.Sleep(1 * time.Second)
// 		}
// 	}()

// 	defer registry.Deregister(ctx, instanceID, serviceName)

// 	metadataGateway := mdgateway.New(registry)
// 	ctrl := users.New(metadataGateway)
// 	h := httphandler.New(ctrl)

// 	//http.Handle("/users", http.HandleFunc(h.GetOrgDetails))
// 	http.HandleFunc("/users", h.GetUserDetails)
// 	if err := http.ListenAndServe(":8083", nil); err != nil {
// 		panic(err)
// 	}
// }
