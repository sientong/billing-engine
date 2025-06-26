package repository

import (
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type BillingScheduleRepo interface {
	CreateBillingSchedule(ctx context.Context, tx *sql.Tx, billingSchedule domain.BillingSchedule) (domain.BillingSchedule, error)
	UpdateBillingSchedule(ctx context.Context, tx *sql.Tx, billingSchedule domain.BillingSchedule) (domain.BillingSchedule, error)
	GetBillingScheduleByUserId(ctx context.Context, tx *sql.Tx, userId string) ([]domain.BillingSchedule, error)
	GetBillingScheduleByLoanId(ctx context.Context, tx *sql.Tx, loanId string) ([]domain.BillingSchedule, error)
}
