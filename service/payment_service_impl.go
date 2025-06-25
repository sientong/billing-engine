package service

import (
	"billing-engine/helper"
	"billing-engine/model/domain"
	"billing-engine/protobuff/pb"
	"billing-engine/repository"
	"context"
	"database/sql"
)

type PaymentServiceImpl struct {
	pb.UnimplementedPaymentServiceServer
	repo repository.PaymentRepository
	repository.LoanRepository
	repository.BillingScheduleRepo
	DB *sql.DB
}

func NewPaymentService(repo repository.PaymentRepository, db *sql.DB) *PaymentServiceImpl {
	return &PaymentServiceImpl{repo: repo, DB: db}
}

func (payment *PaymentServiceImpl) MakePayment(ctx context.Context, req *pb.MakePaymentRequest) (*pb.MakePaymentResponse, error) {

	tx, err := payment.DB.Begin()
	helper.PanicIfError(err)
	defer helper.CommitOrRollback(tx)

	loan := payment.FindLoanById(ctx, tx, req.LoanId)
	if loan.ID == "" {
		return &pb.MakePaymentResponse{}, nil
	}

	unpaidBillingSchedule := payment.GetBillingScheduleByLoanId(ctx, tx, req.LoanId)

	if len(unpaidBillingSchedule) == 0 {
		return &pb.MakePaymentResponse{}, nil
	}

	var storedPayment domain.Payment

	// paid billing schedule based on paid amount
	for _, billingSchedule := range unpaidBillingSchedule {

		// keep on paying until amount is less than or equal to amount due
		if billingSchedule.AmountDue <= req.Amount {
			req.Amount -= billingSchedule.AmountDue
			billingSchedule.Status = "paid"

			newPayment := domain.Payment{
				LoanID:            billingSchedule.LoanID,
				BillingScheduleID: billingSchedule.ID,
				Amount:            billingSchedule.AmountDue,
				PaymentMethod:     req.PaymentMethod,
				Status:            "paid",
			}

			payment.UpdateBillingSchedule(ctx, tx, billingSchedule)
			storedPayment = payment.repo.MakePayment(ctx, tx, newPayment)

			loan.OutstandingBalance -= billingSchedule.AmountDue
			loan.TotalPayment += billingSchedule.AmountDue

			payment.UpdateLoan(ctx, tx, loan)
		}
	}

	return &pb.MakePaymentResponse{
		PaymentId:   storedPayment.ID,
		LoanId:      storedPayment.LoanID,
		Amount:      storedPayment.Amount,
		Status:      storedPayment.Status,
		PaymentDate: storedPayment.PaymentDate}, nil
}
