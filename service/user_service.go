package service

import (
	pb "billing-engine/protobuff/pb"
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error)
	UpdateDeliquentStatus(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error)
	IsDelinquent(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error)
}
