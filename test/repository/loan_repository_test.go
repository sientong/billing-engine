package repository_test

import (
	"billing-engine/model/domain"
	"billing-engine/repository"
	"billing-engine/test"
	"context"
	"testing"
)

var newLoan = domain.Loan{Amount: 1000000.0, InterestRate: 0.05, InterestAmount: 50000.0, TermMonths: 12, TotalPayment: 1050000.0, OutstandingBalance: 0.0, Status: "active"}

func TestLoanRepository_CreateLoan(t *testing.T) {

	db := test.SetupTestDB()
	defer db.Close()

	user, err := test.CreateUser(db)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	loan, err := test.CreateLoan(db, user.ID, newLoan.Amount, newLoan.InterestRate, newLoan.InterestAmount, newLoan.TermMonths, newLoan.OutstandingBalance, newLoan.TotalPayment)
	if err != nil {
		t.Fatalf("Failed to create loan: %v", err)
	}

	if loan.ID == "" {
		t.Error("Expected loan ID to be set, got empty string")
	}

	if loan.CreatedAt == "" {
		t.Error("Expected loan CreatedAt to be set, got empty string")
	}

	if loan.UpdatedAt == "" {
		t.Error("Expected loan UpdatedAt to be set, got empty string")
	}

	if loan.UserID != user.ID {
		t.Errorf("Expected loan UserID to be %s, got %s", user.ID, loan.UserID)
	}

	if loan.Amount != 1000000.0 {
		t.Errorf("Expected loan Amount to be 1000000.0, got %f", loan.Amount)
	}

	if loan.InterestRate != 0.05 {
		t.Errorf("Expected loan InterestRate to be 0.05, got %f", loan.InterestRate)
	}

	if loan.TermMonths != 12 {
		t.Errorf("Expected loan TermMonths to be 12, got %d", loan.TermMonths)
	}

	if loan.TotalPayment != 0.0 {
		t.Errorf("Expected loan TotalPayment to be 0.0, got %f", loan.TotalPayment)
	}

	if loan.OutstandingBalance != 1050000.0 {
		t.Errorf("Expected loan OutstandingBalance to be 1050000.0, got %f", loan.OutstandingBalance)
	}

	if loan.Status != "active" {
		t.Errorf("Expected loan Status to be 'active', got '%s'", loan.Status)
	}

	test.TruncateAllTables(db)
}

func TestLoanRepository_UpdateLoan(t *testing.T) {

	db := test.SetupTestDB()
	defer db.Close()

	user, err := test.CreateUser(db)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	loan, err := test.CreateLoan(db, user.ID, newLoan.Amount, newLoan.InterestRate, newLoan.InterestAmount, newLoan.TermMonths, newLoan.OutstandingBalance, newLoan.TotalPayment)
	if err != nil {
		t.Fatalf("Failed to create loan: %v", err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	loan.TotalPayment = 105000.0
	loan.OutstandingBalance = 900000.0
	loan = repository.NewLoanRepository().UpdateLoan(ctx, tx, loan)

	if loan.TotalPayment != 105000.0 {
		t.Errorf("Expected loan TotalPayment to be 105000.0, got %f", loan.TotalPayment)
	}

	if loan.OutstandingBalance != 900000.0 {
		t.Errorf("Expected loan OutstandingBalance to be 900000.0, got %f", loan.OutstandingBalance)
	}

}
