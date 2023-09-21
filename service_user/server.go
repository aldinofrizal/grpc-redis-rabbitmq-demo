package main

import (
	"context"
	"example/pb"

	"google.golang.org/protobuf/types/known/emptypb"
)

var Users []*pb.User

type UsersServer struct {
	pb.UnimplementedUsersServer
}

func NewUsersServer() UsersServer {
	return UsersServer{}
}

func (u UsersServer) AddUser(ctx context.Context, param *pb.User) (*pb.User, error) {
	Users = append(Users, param)
	return &pb.User{}, nil
}

func (u UsersServer) GetUsers(ctx context.Context, param *emptypb.Empty) (*pb.ListUser, error) {
	result := &pb.ListUser{
		List: Users,
	}
	return result, nil
}
