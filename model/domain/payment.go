package domain

type Payment struct {
	ID                string  `json:"id"`
	LoanID            string  `json:"loan_id"`
	UserID            string  `json:"user_id"`
	BillingScheduleID string  `json:"billing_schedule_id"` // ID of the associated billing schedule
	Amount            float64 `json:"amount"`
	PaymentDate       string  `json:"payment_date"`   // Date when the payment was made
	PaymentMethod     string  `json:"payment_method"` // e.g., "credit_card", "bank_transfer"
	Status            string  `json:"payment_status"` // e.g., "completed", "pending", "failed"
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
}
