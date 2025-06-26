package service

import (
	"billing-engine/helper"
	"billing-engine/model/domain"
	pb "billing-engine/protobuff/pb"
	"billing-engine/repository"
	"context"
	"database/sql"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type LoanServiceImpl struct {
	pb.UnimplementedLoanServiceServer
	repo        repository.LoanRepository
	billingRepo repository.BillingScheduleRepo
	DB          *sql.DB
}

func NewLoanService(repo repository.LoanRepository, billingRepo repository.BillingScheduleRepo, db *sql.DB) *LoanServiceImpl {
	return &LoanServiceImpl{repo: repo, billingRepo: billingRepo, DB: db}
}

func (loan *LoanServiceImpl) CreateNewLoan(ctx context.Context, req *pb.CreateNewLoanRequest) (*pb.LoanResponse, error) {

	interestAmount := req.Amount * req.InterestRate / 100
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
	helper.CheckErrorOrReturn(err)
	defer helper.CommitOrRollback(tx)

	createdLoan, err := loan.repo.CreateLoan(ctx, tx, newLoan)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create new loan: %v", err)
	}

	now := time.Now()
	paymentDueDate := now.AddDate(0, 0, 7)
	paymentDueDateStr := paymentDueDate.Format("2006-01-02 15:04:05.000000-07") // weekly due date
	amountDue := outstandingBalance / float64(req.TermMonths)                   // split amount due date equally

	for range int(req.TermMonths) {
		newBillingSchedule := domain.BillingSchedule{
			LoanID:         createdLoan.ID,
			PaymentDueDate: paymentDueDateStr,
			AmountDue:      amountDue,
			Status:         "pending",
		}

		_, err := loan.billingRepo.CreateBillingSchedule(ctx, tx, newBillingSchedule)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to create new billing schedules: %v", err)
		}

		paymentDueDate = paymentDueDate.AddDate(0, 0, 7)
		paymentDueDateStr = paymentDueDate.Format("2006-01-02 15:04:05.000000-07") // weekly due date
	}

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
	helper.CheckErrorOrReturn(err)
	defer helper.CommitOrRollback(tx)

	foundLoan, err := loan.repo.FindLoanById(ctx, tx, req.LoanId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get loan by id: %v", err)
	}
	return &pb.OutstandingResponse{
		LoanId:             foundLoan.ID,
		UserId:             foundLoan.UserID,
		OutstandingBalance: foundLoan.OutstandingBalance,
	}, nil
}
