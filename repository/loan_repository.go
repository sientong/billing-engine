package repository

import (
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type LoanRepository interface {
	CreateLoan(ctx context.Context, tx *sql.Tx, loan domain.Loan) domain.Loan
	UpdateLoan(ctx context.Context, tx *sql.Tx, loan domain.Loan) domain.Loan
	FindLoanById(ctx context.Context, tx *sql.Tx, loanId string) domain.Loan
}
