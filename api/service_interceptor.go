package main

import (
	"context"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type UserServiceInterceptor struct{}

func NewUserServiceInterceptor() UserServiceInterceptor {
	return UserServiceInterceptor{}
}

func (interceptor UserServiceInterceptor) EmbedAuthCredentials() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		authenticatedCtx := metadata.AppendToOutgoingContext(ctx, "app_authorize_id", os.Getenv("APP_AUTHORIZE_ID"))
		return invoker(authenticatedCtx, method, req, reply, cc, opts...)
	}
}
