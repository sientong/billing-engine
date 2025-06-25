package service

import (
	"billing-engine/helper"
	"billing-engine/model/domain"
	pb "billing-engine/protobuff/pb"
	"billing-engine/repository"
	"context"
	"database/sql"
)

type UserServiceImpl struct {
	pb.UnimplementedUserServiceServer
	repo repository.UserRepository
	DB   *sql.DB
}

func NewUserService(repo repository.UserRepository, db *sql.DB) *UserServiceImpl {
	return &UserServiceImpl{repo: repo, DB: db}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *pb.CreateUserRequest) *pb.UserResponse {

	user := domain.User{
		Name:           req.Name,
		IdentityNumber: req.IdentityNumber,
		IsDelinquent:   false,
		IsActive:       true,
	}

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	createdUser := s.repo.Create(ctx, tx, user)

	return &pb.UserResponse{
		Id:             createdUser.ID,
		Name:           createdUser.Name,
		IdentityNumber: createdUser.IdentityNumber,
		IsDelinquent:   createdUser.IsDelinquent,
		IsActive:       createdUser.IsActive,
		CreatedAt:      createdUser.CreatedAt,
		UpdatedAt:      createdUser.UpdatedAt,
	}
}

func (s *UserServiceImpl) UpdateDeliquentStatus(ctx context.Context, req *pb.UpdateDeliquentStatusRequest) *pb.UserResponse {
	user := domain.User{
		IdentityNumber: req.IdentityNumber,
		IsDelinquent:   req.IsDelinquent,
	}

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	updatedUser := s.repo.Update(ctx, tx, user)
	return &pb.UserResponse{
		Id:             updatedUser.ID,
		Name:           updatedUser.Name,
		IdentityNumber: updatedUser.IdentityNumber,
		IsDelinquent:   updatedUser.IsDelinquent,
		IsActive:       updatedUser.IsActive,
		CreatedAt:      updatedUser.CreatedAt,
		UpdatedAt:      updatedUser.UpdatedAt,
	}
}

func (s *UserServiceImpl) IsDelinquent(ctx context.Context, req *pb.GetDeliquencyStatusRequest) *pb.DeliquencyStatusResponse {

	user := domain.User{
		IdentityNumber: req.IdentityNumber,
	}

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	foundUser, err := s.repo.FindByIdentityNumber(ctx, nil, user.IdentityNumber)
	if err != nil {
		helper.PanicIfError(err)
	}

	return &pb.DeliquencyStatusResponse{
		IdentityNumber: foundUser.IdentityNumber,
		IsDelinquent:   foundUser.IsDelinquent,
	}
}
