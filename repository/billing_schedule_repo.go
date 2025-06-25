package repository

import (
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type BillingScheduleRepo interface {
	CreateBillingSchedule(ctx context.Context, tx *sql.Tx, billingSchedule domain.BillingSchedule) domain.BillingSchedule
	UpdateBillingSchedule(ctx context.Context, tx *sql.Tx, billingSchedule domain.BillingSchedule) domain.BillingSchedule
	GetBillingScheduleByUserId(ctx context.Context, tx *sql.Tx, userId string) []domain.BillingSchedule
	GetBillingScheduleByLoanId(ctx context.Context, tx *sql.Tx, loanId string) []domain.BillingSchedule
}
