package repository

import (
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type PaymentRepository interface {
	MakePayment(ctx context.Context, tx *sql.Tx, payment domain.Payment) domain.Payment
	UpdatePayment(ctx context.Context, tx *sql.Tx, payment domain.Payment) domain.Payment
	GetPaymentsByLoanId(ctx context.Context, tx *sql.Tx, paymentId int) []domain.Payment
	GetPaymentsByUserId(ctx context.Context, tx *sql.Tx, userId int) []domain.Payment
	GetPaymentById(ctx context.Context, tx *sql.Tx, paymentId int) domain.Payment
}
