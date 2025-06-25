package service

import (
	pb "billing-engine/protobuff/pb"
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, req *pb.CreateUserRequest) *pb.UserResponse
	// UpdateDeliquentStatus(ctx context.Context, req *pb.CreateUserRequest) *pb.UserResponse
	// IsDelinquent(ctx context.Context, req *pb.CreateUserRequest) *pb.UserResponse
}
