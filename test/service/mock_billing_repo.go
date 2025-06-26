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

func (b *MockBillingScheduleRepo) CreateBillingSchedule(ctx context.Context, tx *sql.Tx, billingSchedule domain.BillingSchedule) (domain.BillingSchedule, error) {
	b.CalledWith = billingSchedule
	return b.ReturnBillingSchedule, nil
}

func (b *MockBillingScheduleRepo) UpdateBillingSchedule(ctx context.Context, tx *sql.Tx, billingSchedule domain.BillingSchedule) (domain.BillingSchedule, error) {
	b.CalledWith = billingSchedule
	return b.ReturnBillingSchedule, nil
}

func (b *MockBillingScheduleRepo) GetBillingScheduleByUserId(ctx context.Context, tx *sql.Tx, userId string) ([]domain.BillingSchedule, error) {
	return b.ReturnBillingSchedules, nil
}

func (b *MockBillingScheduleRepo) GetBillingScheduleByLoanId(ctx context.Context, tx *sql.Tx, loanId string) ([]domain.BillingSchedule, error) {
	return b.ReturnBillingSchedules, nil
}
