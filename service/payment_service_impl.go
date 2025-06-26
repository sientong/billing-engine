package service

import (
	"billing-engine/helper"
	"billing-engine/model/domain"
	"billing-engine/protobuff/pb"
	"billing-engine/repository"
	"context"
	"database/sql"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentServiceImpl struct {
	pb.UnimplementedPaymentServiceServer
	repo        repository.PaymentRepository
	loanRepo    repository.LoanRepository
	billingRepo repository.BillingScheduleRepo
	DB          *sql.DB
}

func NewPaymentService(repo repository.PaymentRepository, loanRepo repository.LoanRepository, billingRepo repository.BillingScheduleRepo, db *sql.DB) *PaymentServiceImpl {
	return &PaymentServiceImpl{repo: repo, loanRepo: loanRepo, billingRepo: billingRepo, DB: db}
}

func (payment *PaymentServiceImpl) MakePayment(ctx context.Context, req *pb.MakePaymentRequest) (*pb.MakePaymentResponse, error) {

	tx, err := payment.DB.Begin()
	helper.CheckErrorOrReturn(err)
	defer helper.CommitOrRollback(tx)

	loan, err := payment.loanRepo.FindLoanById(ctx, tx, req.LoanId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get loan: %v", err)
	}

	if loan.ID == "" {
		log.Println("Loan not found for loan ID:", req.LoanId)
		return &pb.MakePaymentResponse{}, nil
	}

	billingSchedules, err := payment.billingRepo.GetBillingScheduleByLoanId(ctx, tx, req.LoanId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get billing schedule: %v", err)
	}

	if len(billingSchedules) == 0 {
		log.Println("No unpaid billing schedule found")
		return &pb.MakePaymentResponse{}, nil
	}

	var storedPayment domain.Payment
	now := time.Now()
	nowStr := now.Format("2006-01-02 15:04:05.000000-07")

	// paid billing schedule based on paid amount
	var numberOfPaidBillingSchedule int
	var minimumAmount float64
	for _, billingSchedule := range billingSchedules {

		minimumAmount = billingSchedule.AmountDue
		log.Println("Billing schedule:", billingSchedule.ID, "Amount due:", billingSchedule.AmountDue, "Status:", billingSchedule.Status)

		// skip paid billing schedule
		if billingSchedule.Status == "paid" {
			continue
		}

		// keep on paying until amount is less than or equal to amount due
		if billingSchedule.AmountDue <= req.Amount {
			req.Amount -= billingSchedule.AmountDue
			billingSchedule.Status = "paid"

			newPayment := domain.Payment{
				LoanID:            billingSchedule.LoanID,
				BillingScheduleID: billingSchedule.ID,
				UserID:            loan.UserID,
				Amount:            billingSchedule.AmountDue,
				PaymentMethod:     req.PaymentMethod,
				PaymentDate:       nowStr,
				Status:            "completed",
			}

			paidBillingSchedule, err := payment.billingRepo.UpdateBillingSchedule(ctx, tx, billingSchedule)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "Failed to update billing schedule transaction: %v", err)
			}

			log.Println("Paid billing schedule:", paidBillingSchedule.ID, "Amount due:", paidBillingSchedule.AmountDue, "Status:", paidBillingSchedule.Status)
			storedPayment, err = payment.repo.MakePayment(ctx, tx, newPayment)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "Failed to store payment transaction: %v", err)
			}

			loan.OutstandingBalance -= billingSchedule.AmountDue
			loan.TotalPayment += billingSchedule.AmountDue

			payment.loanRepo.UpdateLoan(ctx, tx, loan)
			numberOfPaidBillingSchedule++
		}
	}

	log.Println("Stored payment:", storedPayment.ID, "Loan ID:", storedPayment.LoanID, "Amount:", storedPayment.Amount, "Status:", storedPayment.Status, "Payment date:", storedPayment.PaymentDate)

	if numberOfPaidBillingSchedule == 0 {
		return &pb.MakePaymentResponse{}, status.Errorf(codes.Internal, "Insufficient amount to pay: %f, require minimum amount: %f", req.Amount, minimumAmount)
	}

	return &pb.MakePaymentResponse{
		PaymentId:   storedPayment.ID,
		LoanId:      storedPayment.LoanID,
		Amount:      storedPayment.Amount,
		Status:      storedPayment.Status,
		PaymentDate: storedPayment.PaymentDate}, nil
}
