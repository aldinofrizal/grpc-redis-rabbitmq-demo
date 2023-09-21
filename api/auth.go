package main

import (
	"context"
	"example/pb"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Authenticator struct {
	UserService pb.UsersClient
}

func NewAuthenticator(s pb.UsersClient) Authenticator {
	return Authenticator{UserService: s}
}

func (auth Authenticator) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token := ctx.Request().Header.Get("authorization")

		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, ErrorResponse{"invalid token"})
		}
		user, err := auth.UserService.VerifyToken(context.Background(), &pb.Token{Token: token})

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, ErrorResponse{"invalid token"})
		}

		ctx.Set("authUser", user)
		return next(ctx)
	}
}
