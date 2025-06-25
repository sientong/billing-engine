package repository_test

import (
	"billing-engine/model/domain"
	"billing-engine/repository"
	"billing-engine/test"
	"context"
	"testing"
)

var newBillingSchedule = domain.BillingSchedule{PaymentDueDate: "2023-06-30", AmountDue: 100000.0, Status: "pending"}
var newSecondBillingSchedule = domain.BillingSchedule{PaymentDueDate: "2023-07-30", AmountDue: 100000.0, Status: "pending"}

func TestBillingSchedule_CreateBillingSchedule(t *testing.T) {

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

	billingSchedule, err := test.CreateBillingSchedule(db, user.ID, loan.ID, newBillingSchedule.AmountDue, newBillingSchedule.PaymentDueDate, newBillingSchedule.Status)
	if err != nil {
		t.Fatalf("Failed to create billing schedule: %v", err)
	}

	if billingSchedule.ID == "" {
		t.Error("Expected billing schedule ID to be set, got empty string")
	}

	if billingSchedule.CreatedAt == "" {
		t.Error("Expected billing schedule CreatedAt to be set, got empty string")
	}

	if billingSchedule.UpdatedAt == "" {
		t.Error("Expected billing schedule UpdatedAt to be set, got empty string")
	}

	if billingSchedule.LoanID != loan.ID {
		t.Errorf("Expected billing schedule LoanID to be %s, got %s", loan.ID, billingSchedule.LoanID)
	}

	if billingSchedule.AmountDue != newBillingSchedule.AmountDue {
		t.Errorf("Expected billing schedule AmountDue to be %f, got %f", newBillingSchedule.AmountDue, billingSchedule.AmountDue)
	}

	if billingSchedule.PaymentDueDate != newBillingSchedule.PaymentDueDate {
		t.Errorf("Expected billing schedule PaymentDueDate to be %s, got %s", newBillingSchedule.PaymentDueDate, billingSchedule.PaymentDueDate)
	}

	if billingSchedule.Status != newBillingSchedule.Status {
		t.Errorf("Expected billing schedule Status to be %s, got %s", newBillingSchedule.Status, billingSchedule.Status)
	}

	test.TruncateAllTables(db)
}

func TestBillingSchedule_UpdateBillingSchedule(t *testing.T) {

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

	billingSchedule, err := test.CreateBillingSchedule(db, user.ID, loan.ID, newBillingSchedule.AmountDue, newBillingSchedule.PaymentDueDate, newBillingSchedule.Status)
	if err != nil {
		t.Fatalf("Failed to create billing schedule: %v", err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	billingSchedule.Status = "paid"
	billingSchedule.PaymentDueDate = "2023-06-30"

	billingSchedule = repository.NewBillingScheduleRepo().UpdateBillingSchedule(ctx, tx, billingSchedule)

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	if billingSchedule.AmountDue != 100000.0 {
		t.Errorf("Expected billing schedule AmountDue to be 100000.0, got %f", billingSchedule.AmountDue)
	}

	if billingSchedule.Status != "paid" {
		t.Errorf("Expected billing schedule Status to be 'paid', got '%s'", billingSchedule.Status)
	}

	if billingSchedule.PaymentDueDate != "2023-06-30" {
		t.Errorf("Expected billing schedule PaymentDueDate to be '2023-06-30', got '%s'", billingSchedule.PaymentDueDate)
	}

	test.TruncateAllTables(db)
}

func TestBillingSchedule_GetBillingScheduleByUserId(t *testing.T) {

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

	billingSchedule, err := test.CreateBillingSchedule(db, user.ID, loan.ID, newBillingSchedule.AmountDue, newBillingSchedule.PaymentDueDate, newBillingSchedule.Status)
	if err != nil {
		t.Fatalf("Failed to create billing schedule: %v", err)
	}

	secondBillingSchedule, err := test.CreateBillingSchedule(db, user.ID, loan.ID, newSecondBillingSchedule.AmountDue, newSecondBillingSchedule.PaymentDueDate, newSecondBillingSchedule.Status)
	if err != nil {
		t.Fatalf("Failed to create second billing schedule: %v", err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	billingSchedules := repository.NewBillingScheduleRepo().GetBillingScheduleByUserId(ctx, tx, user.ID)

	t.Logf(billingSchedules[0].CreatedAt)

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	if len(billingSchedules) != 2 {
		t.Errorf("Expected 2 billing schedules, got %d", len(billingSchedules))
	}

	if billingSchedules[0].ID != billingSchedule.ID {
		t.Errorf("Expected billing schedule ID to be %s, got %s", billingSchedule.ID, billingSchedules[0].ID)
	}

	if billingSchedules[1].ID != secondBillingSchedule.ID {
		t.Errorf("Expected second billing schedule ID to be %s, got %s", secondBillingSchedule.ID, billingSchedules[1].ID)
	}

	if billingSchedules[0].CreatedAt == billingSchedule.CreatedAt {
		t.Errorf("Expected billing createdat %s", billingSchedule.CreatedAt)
	}

	test.TruncateAllTables(db)
}

func TestBillingSchedule_GetBillingScheduleByLoanId(t *testing.T) {

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

	billingSchedule, err := test.CreateBillingSchedule(db, user.ID, loan.ID, newBillingSchedule.AmountDue, newBillingSchedule.PaymentDueDate, newBillingSchedule.Status)
	if err != nil {
		t.Fatalf("Failed to create billing schedule: %v", err)
	}

	secondBillingSchedule, err := test.CreateBillingSchedule(db, user.ID, loan.ID, newSecondBillingSchedule.AmountDue, newSecondBillingSchedule.PaymentDueDate, newSecondBillingSchedule.Status)
	if err != nil {
		t.Fatalf("Failed to create second billing schedule: %v", err)
	}

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	billingSchedules := repository.NewBillingScheduleRepo().GetBillingScheduleByLoanId(ctx, tx, loan.ID)

	err = tx.Commit()
	if err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	if len(billingSchedules) != 2 {
		t.Errorf("Expected 2 billing schedules, got %d", len(billingSchedules))
	}

	if billingSchedules[0].ID != billingSchedule.ID {
		t.Errorf("Expected billing schedule ID to be %s, got %s", billingSchedule.ID, billingSchedules[0].ID)
	}

	if billingSchedules[1].ID != secondBillingSchedule.ID {
		t.Errorf("Expected second billing schedule ID to be %s, got %s", secondBillingSchedule.ID, billingSchedules[1].ID)
	}

	test.TruncateAllTables(db)
}
