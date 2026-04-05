package auth

import (
	"context"
	"fmt"
	"log"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(validator *TokenValidator) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		fmt.Printf("AuthInterceptor: %s\n", info.FullMethod)
		log.Println("AuthInterceptor: ", info.FullMethod)
		var publicMethods = map[string]struct{}{
			"/auth.v1.AuthService/Login":    {},
			"/auth.v1.AuthService/Register": {},
		}

		_, ok := publicMethods[info.FullMethod]
		if ok {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "%s: missing metadada", info.FullMethod)
		}
		fmt.Printf("METADATA: %v\n", md)
		log.Println("METADATA: ", md)

		//authHeader := md["authorization"]
		authHeader := md.Get("authorization")
		log.Printf("auth headers: %v\n", authHeader)
		if len(authHeader) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "%s: missing authorization header", info.FullMethod)
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		if len(token) == 0 {
			return nil, status.Errorf(codes.Unauthenticated, "%s: missing token", info.FullMethod)
		}

		authToken, err := validator.Validate(token)
		if err != nil {
			return nil, status.Errorf(
				codes.Unauthenticated,
				"%s: missing metadada: %s\n",
				info.FullMethod, err.Error(),
			)
		}

		ctx = context.WithValue(ctx, "user_id", authToken.UserID)

		return handler(ctx, req)
	}
}
