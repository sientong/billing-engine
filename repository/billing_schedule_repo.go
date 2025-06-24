package repository

import "billing-engine/model/domain"

type BillingScheduleRepo interface {
	CreateBillingSchedule(billingSchedule domain.BillingSchedule) (domain.BillingSchedule, error)
	UpdateBillingSchedule(billingSchedule domain.BillingSchedule) (domain.BillingSchedule, error)
	GetBillingScheduleByUserId(userId int64) ([]domain.BillingSchedule, error)
	GetBillingScheduleByLoanId(loanId int64) ([]domain.BillingSchedule, error)
}
