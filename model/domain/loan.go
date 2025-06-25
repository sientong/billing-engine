package domain

type Loan struct {
	ID                 string  `json:"id"`
	Amount             float64 `json:"amount"`
	InterestRate       float64 `json:"interest_rate"`
	InterestAmount     float64 `json:"interest_amount"`
	TermMonths         int     `json:"term_months"`
	TotalPayment       float64 `json:"total_payment"`
	OutstandingBalance float64 `json:"outstanding_balance"`
	Status             string  `json:"status"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
	UserID             string  `json:"user_id"`
}
