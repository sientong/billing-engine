package service

import (
	"billing-engine/helper"
	"billing-engine/model/domain"
	pb "billing-engine/protobuff/pb"
	"billing-engine/repository"
	"context"
	"database/sql"
	"time"
)

type UserServiceImpl struct {
	pb.UnimplementedUserServiceServer
	repo        repository.UserRepository
	billingRepo repository.BillingScheduleRepo
	DB          *sql.DB
}

func NewUserService(repo repository.UserRepository, billingRepo repository.BillingScheduleRepo, db *sql.DB) *UserServiceImpl {
	return &UserServiceImpl{repo: repo, billingRepo: billingRepo, DB: db}
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserResponse, error) {

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
	}, nil
}

func (s *UserServiceImpl) UpdateDeliquentStatus(ctx context.Context, req *pb.UpdateDeliquentStatusRequest) (*pb.UserResponse, error) {

	tx, err := s.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	user, err := s.repo.FindByIdentityNumber(ctx, tx, req.IdentityNumber)
	if err != nil {
		helper.PanicIfError(err)
	}

	billingSchedules := s.billingRepo.GetBillingScheduleByUserId(ctx, tx, user.ID)

	layout := time.RFC3339Nano

	var unpaidCounter = 0
	for _, billingSchedule := range billingSchedules {
		parsedTime, err := time.Parse(layout, billingSchedule.CreatedAt)
		if err != nil {
			helper.PanicIfError(err)
		}

		now := time.Now()

		// check how many unpaid billing before today
		if parsedTime.Before(now) && billingSchedule.Status == "unpaid" {
			unpaidCounter++
		}
	}

	// don't change delinquent status if unpaid billing is less than two
	if unpaidCounter < 2 {
		return &pb.UserResponse{
			Id:             user.ID,
			Name:           user.Name,
			IdentityNumber: user.IdentityNumber,
			IsDelinquent:   user.IsDelinquent,
			IsActive:       user.IsActive,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
		}, nil
	}

	// change delinquent status to true
	user.IsDelinquent = true
	updatedUser := s.repo.Update(ctx, tx, user)

	return &pb.UserResponse{
		Id:             updatedUser.ID,
		Name:           updatedUser.Name,
		IdentityNumber: updatedUser.IdentityNumber,
		IsDelinquent:   updatedUser.IsDelinquent,
		IsActive:       updatedUser.IsActive,
		CreatedAt:      updatedUser.CreatedAt,
		UpdatedAt:      updatedUser.UpdatedAt,
	}, nil
}

func (s *UserServiceImpl) IsDelinquent(ctx context.Context, req *pb.GetDeliquencyStatusRequest) (*pb.DeliquencyStatusResponse, error) {

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
	}, nil
}
