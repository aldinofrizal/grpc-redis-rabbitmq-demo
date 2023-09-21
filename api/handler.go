package main

import (
	"context"
	"example/pb"
	"net/http"

	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	Service pb.UsersClient
}

func NewUserHandler(pb pb.UsersClient) UserHandler {
	return UserHandler{Service: pb}
}

func (handler UserHandler) Register(ctx echo.Context) error {
	reqBody := UserRequest{}
	if err := ctx.Bind(&reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrorResponse{err.Error()})
	}

	newUser := pb.User{
		Username: reqBody.Username,
		Password: reqBody.Password,
	}

	createdUser, err := handler.Service.AddUser(context.Background(), &newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrorResponse{err.Error()})
	}

	return ctx.JSON(http.StatusCreated, UserResponse{Id: createdUser.Id, Username: createdUser.Username})
}

func (handler UserHandler) FindAll(ctx echo.Context) error {
	users, err := handler.Service.GetUsers(context.Background(), &emptypb.Empty{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, ErrorResponse{err.Error()})
	}

	response := []UserResponse{}
	for _, u := range users.List {
		response = append(response, UserResponse{Id: u.Id, Username: u.Username})
	}

	return ctx.JSON(http.StatusOK, response)
}

func (handler UserHandler) GetToken(ctx echo.Context) error {
	reqBody := UserRequest{}
	if err := ctx.Bind(&reqBody); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, ErrorResponse{err.Error()})
	}

	loginUser := pb.User{
		Username: reqBody.Username,
		Password: reqBody.Password,
	}

	result, err := handler.Service.GetToken(context.Background(), &loginUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, ErrorResponse{err.Error()})
	}

	return ctx.JSON(http.StatusOK, result)
}
