package auth

import (
	"context"
	"fmt"
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
			return nil, status.Error(codes.Unauthenticated, "missing metadada")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		if len(token) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing token")
		}

		authToken, err := validator.Validate(token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "missing metadada: %s\n", err.Error())
		}

		ctx = context.WithValue(ctx, "user_id", authToken.UserID)

		return handler(ctx, req)
	}
}
