package repository

import (
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type LoanRepositoryImpl struct{}

func NewLoanRepository() *LoanRepositoryImpl {
	return &LoanRepositoryImpl{}
}

func (l *LoanRepositoryImpl) CreateLoan(ctx context.Context, tx *sql.Tx, loan domain.Loan) (domain.Loan, error) {
	SQL := "INSERT INTO loans (user_id, amount, interest_rate, term, status) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at"

	err := tx.QueryRowContext(ctx, SQL, loan.UserID, loan.Amount, loan.InterestRate, loan.TermMonths, "active").Scan(&loan.ID, &loan.CreatedAt, &loan.UpdatedAt)
	if err != nil {
		return domain.Loan{}, err
	}
	return loan, nil
}

func (l *LoanRepositoryImpl) UpdateLoan(ctx context.Context, tx *sql.Tx, loan domain.Loan) (domain.Loan, error) {

	SQL := "UPDATE loans SET total_payment = $1, outstanding_balance = $2, updated_at = NOW() WHERE id = $5 RETURNING updated_at"
	err := tx.QueryRowContext(ctx, SQL, loan.TotalPayment, loan.OutstandingBalance, loan.TermMonths, loan.ID).Scan(&loan.UpdatedAt)
	if err != nil {
		return domain.Loan{}, err
	}

	return loan, nil
}
