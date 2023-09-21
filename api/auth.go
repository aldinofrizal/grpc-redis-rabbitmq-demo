package main

import (
	"context"
	"encoding/json"
	"example/pb"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Authenticator struct {
	UserService pb.UsersClient
	Cache       CacheStorage
}

func NewAuthenticator(s pb.UsersClient, c CacheStorage) Authenticator {
	return Authenticator{UserService: s, Cache: c}
}

func (auth Authenticator) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		token := ctx.Request().Header.Get("authorization")

		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, ErrorResponse{"invalid token"})
		}

		user, err := auth.GetUserCache(token)

		if err != nil {
			user, err = auth.UserService.VerifyToken(context.Background(), &pb.Token{Token: token})
			auth.SaveUserCache(token, user)
		}

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, ErrorResponse{"invalid token"})
		}

		ctx.Set("authUser", user)
		return next(ctx)
	}
}

func (auth Authenticator) SaveUserCache(token string, user *pb.User) {
	auth.Cache.Set(context.Background(), token, auth.MarshallUser(user))
}

func (auth Authenticator) GetUserCache(token string) (*pb.User, error) {
	val, err := auth.Cache.Get(context.Background(), token)
	user := &pb.User{}
	if err != nil {
		return user, err
	}

	err = auth.UnmarhsallUser(val, user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (auth Authenticator) MarshallUser(u *pb.User) string {
	return fmt.Sprintf(`{
		"id": "%s",
		"username": "%s"
	}`, u.Id, u.Username)
}

func (auth Authenticator) UnmarhsallUser(v string, user *pb.User) error {
	err := json.Unmarshal([]byte(v), user)
	if err != nil {
		return err
	}

	return nil
}
