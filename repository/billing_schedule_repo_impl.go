package repository

import (
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type BillingScheduleRepoImpl struct{}

func NewBillingScheduleRepo() *BillingScheduleRepoImpl {
	return &BillingScheduleRepoImpl{}
}

func (b *BillingScheduleRepoImpl) CreateBillingSchedule(ctx context.Context, tx *sql.Tx, billingSchedule domain.BillingSchedule) (domain.BillingSchedule, error) {

	SQL := `INSERT INTO billing_schedules (user_id, loan_id, payment_amount, payment_due_date, payment_status)
			VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`

	err := tx.QueryRowContext(ctx, SQL, billingSchedule.UserID, billingSchedule.LoanID, billingSchedule.PaymentAmount, billingSchedule.PaymentDueDate, billingSchedule.PaymentStatus).
		Scan(&billingSchedule.ID, &billingSchedule.CreatedAt, &billingSchedule.UpdatedAt)
	if err != nil {
		return domain.BillingSchedule{}, err
	}
	return billingSchedule, nil
}

func (b *BillingScheduleRepoImpl) UpdateBillingSchedule(ctx context.Context, tx *sql.Tx, billingSchedule domain.BillingSchedule) (domain.BillingSchedule, error) {

	SQL := `UPDATE billing_schedules SET payment_amount = $1, payment_due_date = $2, payment_status = $3, updated_at = NOW()
			WHERE id = $4 RETURNING updated_at`

	err := tx.QueryRowContext(ctx, SQL, billingSchedule.PaymentAmount, billingSchedule.PaymentDueDate, billingSchedule.PaymentStatus, billingSchedule.ID).
		Scan(&billingSchedule.UpdatedAt)
	if err != nil {
		return domain.BillingSchedule{}, err
	}

	return billingSchedule, nil
}

func (b *BillingScheduleRepoImpl) GetBillingScheduleByUserId(ctx context.Context, tx *sql.Tx, userID string) ([]domain.BillingSchedule, error) {
	SQL := `SELECT id, user_id, loan_id, payment_due_date, payment_amount, payment_status, created_at, updated_at
			FROM billing_schedules WHERE user_id = $1`

	rows, err := tx.QueryContext(ctx, SQL, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var billingSchedules []domain.BillingSchedule
	for rows.Next() {
		var billingSchedule domain.BillingSchedule
		err := rows.Scan(&billingSchedule.ID, &billingSchedule.UserID, &billingSchedule.LoanID,
			&billingSchedule.PaymentDueDate, &billingSchedule.PaymentAmount,
			&billingSchedule.PaymentStatus, &billingSchedule.CreatedAt, &billingSchedule.UpdatedAt)
		if err != nil {
			return nil, err
		}
		billingSchedules = append(billingSchedules, billingSchedule)
	}

	return billingSchedules, nil
}

func (b *BillingScheduleRepoImpl) GetBillingScheduleByLoanId(ctx context.Context, tx *sql.Tx, loanID string) ([]domain.BillingSchedule, error) {
	SQL := `SELECT id, user_id, loan_id, payment_due_date, payment_amount, payment_status, created_at, updated_at
			FROM billing_schedules WHERE loan_id = $1`

	rows, err := tx.QueryContext(ctx, SQL, loanID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var billingSchedules []domain.BillingSchedule
	for rows.Next() {
		var billingSchedule domain.BillingSchedule
		err := rows.Scan(&billingSchedule.ID, &billingSchedule.UserID, &billingSchedule.LoanID,
			&billingSchedule.PaymentDueDate, &billingSchedule.PaymentAmount,
			&billingSchedule.PaymentStatus, &billingSchedule.CreatedAt, &billingSchedule.UpdatedAt)
		if err != nil {
			return nil, err
		}
		billingSchedules = append(billingSchedules, billingSchedule)
	}

	return billingSchedules, nil
}
