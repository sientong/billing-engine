package repository

import (
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, tx *sql.Tx, loan domain.Loan) (domain.Loan, error)
	UpdateLoan(ctx context.Context, tx *sql.Tx, loan domain.Loan) (domain.Loan, error)
	FindLoanById(ctx context.Context, tx *sql.Tx, loanId string) (domain.Loan, error)
}
