package repository

import (
	"billing-engine/helper"
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type BillingScheduleRepoImpl struct{}

func NewBillingScheduleRepo() *BillingScheduleRepoImpl {
	return &BillingScheduleRepoImpl{}
}

func (b *BillingScheduleRepoImpl) CreateBillingSchedule(ctx context.Context, tx *sql.Tx, billingSchedule domain.BillingSchedule) (domain.BillingSchedule, error) {

	SQL := `INSERT INTO billing_schedules (loan_id, amount_due, payment_due_date, status)
			VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`

	err := tx.QueryRowContext(ctx, SQL, billingSchedule.LoanID, billingSchedule.AmountDue, billingSchedule.PaymentDueDate, billingSchedule.Status).
		Scan(&billingSchedule.ID, &billingSchedule.CreatedAt, &billingSchedule.UpdatedAt)
	if err != nil {
		helper.CheckErrorOrReturn(err)
		return domain.BillingSchedule{}, err
	}
	return billingSchedule, nil
}

func (b *BillingScheduleRepoImpl) UpdateBillingSchedule(ctx context.Context, tx *sql.Tx, billingSchedule domain.BillingSchedule) (domain.BillingSchedule, error) {

	SQL := `UPDATE billing_schedules SET amount_due = $1, payment_due_date = $2, status = $3, updated_at = NOW()
			WHERE id = $4 RETURNING updated_at`

	err := tx.QueryRowContext(ctx, SQL, billingSchedule.AmountDue, billingSchedule.PaymentDueDate, billingSchedule.Status, billingSchedule.ID).
		Scan(&billingSchedule.UpdatedAt)
	if err != nil {
		helper.CheckErrorOrReturn(err)
		return domain.BillingSchedule{}, err
	}

	return billingSchedule, nil
}

func (b *BillingScheduleRepoImpl) GetBillingScheduleByUserId(ctx context.Context, tx *sql.Tx, userID string) ([]domain.BillingSchedule, error) {
	SQL := `SELECT bs.id, bs.loan_id, bs.payment_due_date, bs.amount_due, bs.status, bs.created_at, bs.updated_at
			FROM billing_schedules bs JOIN loans l ON bs.loan_id = l.id
			WHERE user_id = $1`

	rows, err := tx.QueryContext(ctx, SQL, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []domain.BillingSchedule{}, nil
		}
		helper.CheckErrorOrReturn(err)
		return nil, err
	}
	defer rows.Close()

	var billingSchedules []domain.BillingSchedule
	for rows.Next() {
		var billingSchedule domain.BillingSchedule
		err := rows.Scan(&billingSchedule.ID, &billingSchedule.LoanID,
			&billingSchedule.PaymentDueDate, &billingSchedule.AmountDue,
			&billingSchedule.Status, &billingSchedule.CreatedAt, &billingSchedule.UpdatedAt)

		if err != nil {
			helper.CheckErrorOrReturn(err)
			return nil, err
		}
		billingSchedules = append(billingSchedules, billingSchedule)
	}

	return billingSchedules, nil
}

func (b *BillingScheduleRepoImpl) GetBillingScheduleByLoanId(ctx context.Context, tx *sql.Tx, loanID string) ([]domain.BillingSchedule, error) {
	SQL := `SELECT id, loan_id, payment_due_date, amount_due, status, created_at, updated_at
			FROM billing_schedules WHERE loan_id = $1`

	rows, err := tx.QueryContext(ctx, SQL, loanID)
	if err != nil {
		if err == sql.ErrNoRows {
			return []domain.BillingSchedule{}, err
		}
		helper.CheckErrorOrReturn(err)
		return nil, err
	}
	defer rows.Close()

	var billingSchedules []domain.BillingSchedule
	for rows.Next() {
		var billingSchedule domain.BillingSchedule
		err := rows.Scan(&billingSchedule.ID, &billingSchedule.LoanID,
			&billingSchedule.PaymentDueDate, &billingSchedule.AmountDue,
			&billingSchedule.Status, &billingSchedule.CreatedAt, &billingSchedule.UpdatedAt)

		if err != nil {
			helper.CheckErrorOrReturn(err)
			return nil, err
		}

		billingSchedules = append(billingSchedules, billingSchedule)
	}

	return billingSchedules, err
}
