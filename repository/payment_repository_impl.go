package repository

import (
	"billing-engine/helper"
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type PaymentRepositoryImpl struct{}

func NewPaymentRepository() *PaymentRepositoryImpl {
	return &PaymentRepositoryImpl{}
}

func (p *PaymentRepositoryImpl) CreatePayment(ctx context.Context, tx *sql.Tx, payment domain.Payment) (domain.Payment, error) {
	SQL := `INSERT INTO payments (user_id, loan_id, billing_schedule_id, amount, payment_date, payment_method, payment_status) VALUES ($1, $2, $3, $4, %5, $6, $7) RETURNING id, created_at, updated_at`

	err := tx.QueryRowContext(ctx, SQL, payment.UserID, payment.LoanID, payment.BillingScheduleID, payment.Amount, payment.PaymentDate, payment.PaymentMethod, payment.PaymentStatus).
		Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		helper.PanicIfError(err)
	}

	return payment, nil
}

func (p *PaymentRepositoryImpl) UpdatePayment(ctx context.Context, tx *sql.Tx, payment domain.Payment) (domain.Payment, error) {

	SQL := `UPDATE payments SET payment_status = $1, updated_at = NOW() WHERE id = $5 RETURNING updated_at`
	err := tx.QueryRowContext(ctx, SQL, payment.PaymentStatus, payment.ID).
		Scan(&payment.UpdatedAt)
	if err != nil {
		helper.PanicIfError(err)
	}

	return payment, nil
}

func (p *PaymentRepositoryImpl) GetPaymentsByLoanId(ctx context.Context, tx *sql.Tx, loanId int64) ([]domain.Payment, error) {
	var payments []domain.Payment

	SQL := `SELECT id, user_id, loan_id, billing_schedule_id, amount, payment_date, payment_method, payment_status, created_at, updated_at
			FROM payments WHERE loan_id = $1`
	rows, err := tx.QueryContext(ctx, SQL, loanId)
	if err != nil {
		helper.PanicIfError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var payment domain.Payment
		err := rows.Scan(&payment.ID, &payment.UserID, &payment.LoanID, &payment.BillingScheduleID, &payment.Amount, &payment.PaymentDate, &payment.PaymentMethod, &payment.PaymentStatus, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			helper.PanicIfError(err)
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (p *PaymentRepositoryImpl) GetPaymentsByUserId(ctx context.Context, tx *sql.Tx, userId int64) ([]domain.Payment, error) {
	var payments []domain.Payment

	SQL := `SELECT id, user_id, loan_id, billing_schedule_id, amount, payment_date, payment_method, payment_status, created_at, updated_at
			FROM payments WHERE user_id = $1`
	rows, err := tx.QueryContext(ctx, SQL, userId)
	if err != nil {
		helper.PanicIfError(err)
	}
	defer rows.Close()

	for rows.Next() {
		var payment domain.Payment
		err := rows.Scan(&payment.ID, &payment.UserID, &payment.LoanID, &payment.BillingScheduleID, &payment.Amount, &payment.PaymentDate, &payment.PaymentMethod, &payment.PaymentStatus, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			helper.PanicIfError(err)
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (p *PaymentRepositoryImpl) GetPaymentById(ctx context.Context, tx *sql.Tx, paymentId int64) (domain.Payment, error) {
	var payment domain.Payment
	SQL := `SELECT id, user_id, loan_id, billing_schedule_id, amount, payment_date, payment_method, payment_status, created_at, updated_at
			FROM payments WHERE id = $1`
	err := tx.QueryRowContext(ctx, SQL, paymentId).Scan(&payment.ID, &payment.UserID, &payment.LoanID, &payment.BillingScheduleID, &payment.Amount, &payment.PaymentDate, &payment.PaymentMethod, &payment.PaymentStatus, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Payment{}, nil // No payment found
		}
		helper.PanicIfError(err)
	}

	return payment, nil
}
