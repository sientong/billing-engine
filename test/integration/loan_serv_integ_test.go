package integration_test

import (
	"billing-engine/protobuff/pb"
	"billing-engine/test"
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateLoan(userId string) *pb.CreateNewLoanRequest {
	return &pb.CreateNewLoanRequest{
		UserId:       userId,
		Amount:       10000000,
		InterestRate: 10,
		TermMonths:   10,
	}
}

func TestLoanService_Create(t *testing.T) {

	ctx := context.Background()

	req := CreateUser()
	user, err := client.CreateUser(ctx, req)
	require.NoError(t, err)

	secondReq := CreateLoan(user.Id)
	loan, err := loanClient.CreateNewLoan(ctx, secondReq)
	require.NoError(t, err)

	require.Equal(t, user.Id, loan.UserId)
	require.Equal(t, secondReq.Amount, loan.Amount)
	require.Equal(t, secondReq.InterestRate, loan.InterestRate)
	require.Equal(t, secondReq.TermMonths, loan.TermMonths)
	require.NotEmpty(t, loan.LoanId)
	require.NotEmpty(t, loan.CreatedAt)
	require.NotEmpty(t, loan.UpdatedAt)

	test.TruncateAllTables(db)
}

func TestLoanService_GetOutstanding(t *testing.T) {

	ctx := context.Background()

	req := CreateUser()
	user, err := client.CreateUser(ctx, req)
	require.NoError(t, err)

	secondReq := CreateLoan(user.Id)
	loan, err := loanClient.CreateNewLoan(ctx, secondReq)
	require.NoError(t, err)

	thirdReq := &pb.GetOutstandingRequest{LoanId: loan.LoanId}
	getOutstanding, err := loanClient.GetOutstanding(ctx, thirdReq)
	require.NoError(t, err)

	require.Equal(t, loan.LoanId, getOutstanding.LoanId)
	require.Equal(t, user.Id, getOutstanding.UserId)
	require.Equal(t, loan.OutstandingBalance, getOutstanding.OutstandingBalance)

	test.TruncateAllTables(db)
}
