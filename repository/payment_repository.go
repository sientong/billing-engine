package repository

import (
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type PaymentRepository interface {
	MakePayment(ctx context.Context, tx *sql.Tx, payment domain.Payment) (domain.Payment, error)
	UpdatePayment(ctx context.Context, tx *sql.Tx, payment domain.Payment) (domain.Payment, error)
	GetPaymentsByLoanId(ctx context.Context, tx *sql.Tx, paymentId int) ([]domain.Payment, error)
	GetPaymentsByUserId(ctx context.Context, tx *sql.Tx, userId int) ([]domain.Payment, error)
	GetPaymentById(ctx context.Context, tx *sql.Tx, paymentId int) (domain.Payment, error)
}
