package repository

import "billing-engine/model/domain"

type BillingScheduleRepo interface {
	CreateBillingSchedule(billingSchedule domain.BillingSchedule) domain.BillingSchedule
	UpdateBillingSchedule(billingSchedule domain.BillingSchedule) domain.BillingSchedule
	GetBillingScheduleByUserId(userId int64) []domain.BillingSchedule
	GetBillingScheduleByLoanId(loanId int64) []domain.BillingSchedule
}
