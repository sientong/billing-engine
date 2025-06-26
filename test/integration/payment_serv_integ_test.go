package integration_test

import (
	"billing-engine/protobuff/pb"
	"billing-engine/test"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPaymentService_MakePayment(t *testing.T) {
	ctx := context.Background()

	req := CreateUser()
	user, err := client.CreateUser(ctx, req)
	require.NoError(t, err)

	secondReq := CreateLoan(user.Id)
	loan, err := loanClient.CreateNewLoan(ctx, secondReq)
	require.NoError(t, err)

	paymentReq := &pb.MakePaymentRequest{
		LoanId:        loan.LoanId,
		Amount:        1100000,
		PaymentMethod: "bank transfer",
	}

	payment, err := paymentClient.MakePayment(ctx, paymentReq)
	require.NoError(t, err)

	require.Equal(t, loan.LoanId, payment.LoanId)
	require.Equal(t, payment.Amount, paymentReq.Amount)
	require.Equal(t, payment.Status, "completed")
	require.NotEmpty(t, payment.PaymentId)
	require.NotEmpty(t, payment.PaymentDate)

	test.TruncateAllTables(db)
}

func TestPaymentService_MakePaymentWithInsufficientAmount(t *testing.T) {
	ctx := context.Background()

	req := CreateUser()
	user, err := client.CreateUser(ctx, req)
	require.NoError(t, err)

	secondReq := CreateLoan(user.Id)
	loan, err := loanClient.CreateNewLoan(ctx, secondReq)
	require.NoError(t, err)

	paymentReq := &pb.MakePaymentRequest{
		LoanId:        loan.LoanId,
		Amount:        10000,
		PaymentMethod: "bank transfer",
	}

	_, err = paymentClient.MakePayment(ctx, paymentReq)
	require.Error(t, err)
	require.Equal(t, "rpc error: code = Internal desc = Insufficient amount to pay: 10000.000000, require minimum amount: 1100000.000000", err.Error())

	test.TruncateAllTables(db)
}
