package service_test

import (
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type MockBillingScheduleRepo struct {
	CalledWith             domain.BillingSchedule
	ReturnBillingSchedule  domain.BillingSchedule
	ReturnBillingSchedules []domain.BillingSchedule
}

func (b *MockBillingScheduleRepo) CreateBillingSchedule(ctx context.Context, tx *sql.Tx, billingSchedule domain.BillingSchedule) domain.BillingSchedule {
	b.CalledWith = billingSchedule
	return b.ReturnBillingSchedule
}

func (b *MockBillingScheduleRepo) UpdateBillingSchedule(ctx context.Context, tx *sql.Tx, billingSchedule domain.BillingSchedule) domain.BillingSchedule {
	b.CalledWith = billingSchedule
	return b.ReturnBillingSchedule
}

func (b *MockBillingScheduleRepo) GetBillingScheduleByUserId(ctx context.Context, tx *sql.Tx, userId string) []domain.BillingSchedule {
	return []domain.BillingSchedule{}
}

func (b *MockBillingScheduleRepo) GetBillingScheduleByLoanId(ctx context.Context, tx *sql.Tx, loanId string) []domain.BillingSchedule {
	return []domain.BillingSchedule{}
}
