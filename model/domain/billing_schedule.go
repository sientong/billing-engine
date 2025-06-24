package domain

type BillingSchedule struct {
	ID             string  `json:"id"`
	LoanID         string  `json:"loan_id"`
	PaymentDueDate string  `json:"payment_due_date"`
	PaymentAmount  float64 `json:"payment_amount"`
	PaymentStatus  string  `json:"payment_status"` // e.g., "pending", "paid", "overdue"
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
	UserID         int64   `json:"user_id"` // ID of the user associated with the billing schedule
}
