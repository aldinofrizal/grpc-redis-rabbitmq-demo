package main

import (
	"context"
	"errors"
	"example/helpers"
	"example/pb"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/protobuf/types/known/emptypb"
)

var Users []*pb.User

type UsersServer struct {
	pb.UnimplementedUsersServer
	Repository UserRepository
}

func NewUsersServer(rp UserRepository) UsersServer {
	return UsersServer{Repository: rp}
}

func (u UsersServer) AddUser(ctx context.Context, param *pb.User) (*pb.User, error) {
	userPb := pb.User{}
	user, err := u.Repository.CreateUser(param.Username, helpers.HashPassword(param.Password))
	if err != nil {
		return &userPb, errors.New("failed to create user")
	}

	userPb.Id = user.ID.Hex()
	userPb.Username = user.Username
	userPb.Password = user.Password

	return &userPb, nil
}

func (u UsersServer) GetUsers(ctx context.Context, param *emptypb.Empty) (*pb.ListUser, error) {
	result := &pb.ListUser{
		List: []*pb.User{},
	}

	users, err := u.Repository.GetAllUsers()
	if err != nil {
		return result, errors.New("failed to get users")
	}

	for _, u := range users {
		result.List = append(result.List, &pb.User{
			Id:       u.ID.Hex(),
			Username: u.Username,
			Password: u.Password,
		})
	}

	return result, nil
}

func (u UsersServer) GetToken(ctx context.Context, param *pb.User) (*pb.Token, error) {
	result := pb.Token{}
	user, err := u.Repository.GetUserByUsername(param.Username)
	if err != nil {
		return &result, errors.New("invalid credentials")
	}

	validPassword := helpers.ComparePassword(param.Password, user.Password)
	if !validPassword {
		return &result, errors.New("invalid credentials")
	}

	token := helpers.GenerateToken(jwt.MapClaims{"id": user.ID})
	result.Token = token

	return &result, nil
}

func (u UsersServer) VerifyToken(ctx context.Context, param *pb.Token) (*pb.User, error) {
	result := pb.User{}
	user, err := u.Repository.VerifyUserByToken(param.Token)
	if err != nil {
		return &result, errors.New("invalid token")
	}

	result.Id = user.ID.Hex()
	result.Username = user.Username

	return &result, nil
}
