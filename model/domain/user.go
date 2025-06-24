package domain

type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	IdentityNumber string `json:"identity_number"`
	IsDelinquent   bool   `json:"is_delinquent"`
	IsActive       bool   `json:"is_active"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
}
