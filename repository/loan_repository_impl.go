package repository

import (
	"billing-engine/helper"
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type LoanRepositoryImpl struct{}

func NewLoanRepository() *LoanRepositoryImpl {
	return &LoanRepositoryImpl{}
}

func (l *LoanRepositoryImpl) CreateLoan(ctx context.Context, tx *sql.Tx, loan domain.Loan) domain.Loan {
	SQL := "INSERT INTO loans (user_id, amount, interest_rate, term_months, total_payment, outstanding_balance, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at"

	err := tx.QueryRowContext(ctx, SQL, loan.UserID, loan.Amount, loan.InterestRate, loan.TermMonths, loan.TotalPayment, loan.OutstandingBalance, "active").Scan(&loan.ID, &loan.CreatedAt, &loan.UpdatedAt)
	if err != nil {
		helper.PanicIfError(err)
	}

	return loan
}

func (l *LoanRepositoryImpl) UpdateLoan(ctx context.Context, tx *sql.Tx, loan domain.Loan) domain.Loan {

	SQL := "UPDATE loans SET total_payment = $1, outstanding_balance = $2, updated_at = NOW() WHERE id = $3 RETURNING updated_at"
	err := tx.QueryRowContext(ctx, SQL, loan.TotalPayment, loan.OutstandingBalance, loan.ID).Scan(&loan.UpdatedAt)
	if err != nil {
		helper.PanicIfError(err)
	}

	return loan
}
