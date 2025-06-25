package domain

type BillingSchedule struct {
	ID             string  `json:"id"`
	LoanID         string  `json:"loan_id"`
	PaymentDueDate string  `json:"payment_due_date"`
	AmountDue      float64 `json:"amount_due"`
	Status         string  `json:"status"` // e.g., "pending", "paid", "overdue"
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}
