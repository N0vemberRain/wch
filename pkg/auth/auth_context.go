package auth

import (
	"context"

	"google.golang.org/grpc/metadata"
)

func ForwardAuthContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	authHeader := md.Get("authorization")
	if len(authHeader) == 0 {
		return ctx
	}

	return metadata.NewOutgoingContext(
		ctx,
		metadata.Pairs("authorization", authHeader[0]),
	)
}
