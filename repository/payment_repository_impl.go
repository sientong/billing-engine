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

func (p *PaymentRepositoryImpl) MakePayment(ctx context.Context, tx *sql.Tx, payment domain.Payment) (domain.Payment, error) {
	SQL := `INSERT INTO payments (user_id, loan_id, billing_schedule_id, amount, payment_date, payment_method, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id, created_at, updated_at`

	err := tx.QueryRowContext(ctx, SQL, payment.UserID, payment.LoanID, payment.BillingScheduleID, payment.Amount, payment.PaymentDate, payment.PaymentMethod, payment.Status).
		Scan(&payment.ID, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		helper.CheckErrorOrReturn(err)
		return domain.Payment{}, err
	}

	return payment, nil
}

func (p *PaymentRepositoryImpl) UpdatePayment(ctx context.Context, tx *sql.Tx, payment domain.Payment) (domain.Payment, error) {

	SQL := `UPDATE payments SET status = $1, updated_at = NOW() WHERE id = $5 RETURNING updated_at`
	err := tx.QueryRowContext(ctx, SQL, payment.Status, payment.ID).
		Scan(&payment.UpdatedAt)
	if err != nil {
		helper.CheckErrorOrReturn(err)
		return domain.Payment{}, err
	}

	return payment, nil
}

func (p *PaymentRepositoryImpl) GetPaymentsByLoanId(ctx context.Context, tx *sql.Tx, loanId string) ([]domain.Payment, error) {
	var payments []domain.Payment

	SQL := `SELECT id, user_id, loan_id, billing_schedule_id, amount, payment_date, payment_method, payment_status, created_at, updated_at
			FROM payments WHERE loan_id = $1`
	rows, err := tx.QueryContext(ctx, SQL, loanId)
	if err != nil {
		helper.CheckErrorOrReturn(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment domain.Payment
		err := rows.Scan(&payment.ID, &payment.UserID, &payment.LoanID, &payment.BillingScheduleID, &payment.Amount, &payment.PaymentDate, &payment.PaymentMethod, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			helper.CheckErrorOrReturn(err)
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (p *PaymentRepositoryImpl) GetPaymentsByUserId(ctx context.Context, tx *sql.Tx, userId string) ([]domain.Payment, error) {
	var payments []domain.Payment

	SQL := `SELECT id, user_id, loan_id, billing_schedule_id, amount, payment_date, payment_method, status, created_at, updated_at
			FROM payments WHERE user_id = $1`
	rows, err := tx.QueryContext(ctx, SQL, userId)
	if err != nil {
		helper.CheckErrorOrReturn(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var payment domain.Payment
		err := rows.Scan(&payment.ID, &payment.UserID, &payment.LoanID, &payment.BillingScheduleID, &payment.Amount, &payment.PaymentDate, &payment.PaymentMethod, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt)
		if err != nil {
			helper.CheckErrorOrReturn(err)
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (p *PaymentRepositoryImpl) GetPaymentById(ctx context.Context, tx *sql.Tx, paymentId string) (domain.Payment, error) {
	var payment domain.Payment
	SQL := `SELECT id, user_id, loan_id, billing_schedule_id, amount, payment_date, payment_method, status, created_at, updated_at
			FROM payments WHERE id = $1`
	err := tx.QueryRowContext(ctx, SQL, paymentId).Scan(&payment.ID, &payment.UserID, &payment.LoanID, &payment.BillingScheduleID, &payment.Amount, &payment.PaymentDate, &payment.PaymentMethod, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Payment{}, err
		}
		helper.CheckErrorOrReturn(err)
	}

	return payment, nil
}
