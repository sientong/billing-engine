package service

import (
	"billing-engine/helper"
	"billing-engine/model/domain"
	pb "billing-engine/protobuff/pb"
	"billing-engine/repository"
	"context"
	"database/sql"
)

type LoanServiceImpl struct {
	pb.UnimplementedLoanServiceServer
	repo repository.LoanRepository
	DB   *sql.DB
}

func NewLoanService(repo repository.LoanRepository, db *sql.DB) *LoanServiceImpl {
	return &LoanServiceImpl{repo: repo, DB: db}
}

func (loan *LoanServiceImpl) CreateNewLoan(ctx context.Context, req *pb.CreateNewLoanRequest) (*pb.LoanResponse, error) {

	interestAmount := req.Amount * req.InterestRate
	outstandingBalance := req.Amount + interestAmount

	newLoan := domain.Loan{
		Amount:             req.Amount,
		InterestRate:       req.InterestRate,
		InterestAmount:     interestAmount,
		OutstandingBalance: outstandingBalance,
		TermMonths:         int(req.TermMonths),
		TotalPayment:       0,
		Status:             "pending",
		UserID:             req.UserId,
	}

	tx, err := loan.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	createdLoan := loan.repo.CreateLoan(ctx, tx, newLoan)

	return &pb.LoanResponse{
		LoanId:             createdLoan.ID,
		UserId:             createdLoan.UserID,
		Amount:             createdLoan.Amount,
		InterestRate:       createdLoan.InterestRate,
		TermMonths:         int32(createdLoan.TermMonths),
		TotalPayment:       createdLoan.TotalPayment,
		OutstandingBalance: createdLoan.OutstandingBalance,
		Status:             createdLoan.Status,
		CreatedAt:          createdLoan.CreatedAt,
		UpdatedAt:          createdLoan.UpdatedAt,
	}, nil
}

func (loan *LoanServiceImpl) GetOutstanding(ctx context.Context, req *pb.GetOutstandingRequest) (*pb.OutstandingResponse, error) {

	tx, err := loan.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	foundLoan := loan.repo.FindLoanById(ctx, tx, req.LoanId)

	return &pb.OutstandingResponse{
		LoanId:             foundLoan.ID,
		UserId:             foundLoan.UserID,
		OutstandingBalance: foundLoan.OutstandingBalance,
	}, nil
}
