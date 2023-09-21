package main

import (
	"context"
	"errors"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthInterceptor struct{}

func NewAuthInterceptor() AuthInterceptor {
	return AuthInterceptor{}
}

func (auth AuthInterceptor) AuthtenticateApp() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, errors.New("metadata required")
		}

		val := md["app_authorize_id"]
		if len(val) == 0 {
			return nil, errors.New("metadata required")
		}

		if val[0] != os.Getenv("APP_AUTHORIZE_ID") {
			return nil, errors.New("metadata required")
		}

		return handler(ctx, req)
	}
}
