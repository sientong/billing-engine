package test

import (
	"billing-engine/model/domain"
	"billing-engine/repository"
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
)

func ExpectPanicMessage(t *testing.T, fn func(), expectedMessage string) {
	t.Helper()

	defer func() {

		if r := recover(); r != nil {
			var actual string

			switch v := r.(type) {
			case error:
				actual = v.Error()
			case string:
				actual = v
			default:
				actual = fmt.Sprintf("%v", v)
			}

			if actual != expectedMessage {
				t.Errorf("Expected panic message:\n  %q\nbut got:\n  %q", expectedMessage, actual)
			}
		} else {
			t.Errorf("Expected panic with message %q, but no panic occurred", expectedMessage)
		}
	}()

	fn()
}

func CreateUser(db *sql.DB) (domain.User, error) {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
		return domain.User{}, err
	}

	defer tx.Rollback()

	repo := repository.NewUserRepository()
	newUser := domain.User{
		Name:           "John Doe",
		IdentityNumber: "123456789",
		IsActive:       true,
		IsDelinquent:   false,
	}

	result := repo.Create(ctx, tx, newUser)

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
		return domain.User{}, err
	}

	return result, nil
}

func CreateLoan(db *sql.DB, userId string, loanAmount float64, interestRate float64, interestAmount float64, termMonths int, outstandingBalance float64, totalPayment float64) (domain.Loan, error) {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
		return domain.Loan{}, err
	}
	defer tx.Rollback()

	repo := repository.NewLoanRepository()
	newLoan := domain.Loan{
		UserID:             userId,
		Amount:             loanAmount,
		InterestRate:       interestRate,
		InterestAmount:     interestAmount,
		OutstandingBalance: outstandingBalance,
		TotalPayment:       totalPayment,
		TermMonths:         termMonths,
	}

	result := repo.CreateLoan(ctx, tx, newLoan)

	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit transaction: %v", err)
		return domain.Loan{}, err
	}

	return result, nil
}

func CreatePayment(db *sql.DB, userId string, loanId string, billingScheduleId string, amount float64, paymentDate string, paymentMethod string) (domain.Payment, error) {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
		return domain.Payment{}, err
	}

	defer tx.Rollback()

	repo := repository.NewPaymentRepository()
	newPayment := domain.Payment{
		UserID:            userId,
		LoanID:            loanId,
		BillingScheduleID: billingScheduleId,
		Amount:            amount,
		PaymentDate:       paymentDate,
		PaymentMethod:     paymentMethod,
	}

	result := repo.MakePayment(ctx, tx, newPayment)
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit payment transaction: %v", err)
		return domain.Payment{}, err
	}

	return result, nil

}

func CreateBillingSchedule(db *sql.DB, userId string, loanId string, amountDue float64, paymentDueDate string, status string) (domain.BillingSchedule, error) {

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to begin transaction: %v", err)
		return domain.BillingSchedule{}, err
	}

	defer tx.Rollback()

	repo := repository.NewBillingScheduleRepo()
	newBillingSchedule := domain.BillingSchedule{
		LoanID:         loanId,
		AmountDue:      amountDue,
		PaymentDueDate: paymentDueDate,
		Status:         status,
	}

	result := repo.CreateBillingSchedule(ctx, tx, newBillingSchedule)
	if err := tx.Commit(); err != nil {
		log.Fatalf("Failed to commit billing schedule transaction: %v", err)
		return domain.BillingSchedule{}, err
	}

	return result, nil

}
