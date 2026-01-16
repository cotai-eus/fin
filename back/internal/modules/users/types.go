package users

import "time"

// User represents a user in the system
type User struct {
	ID                        string    `json:"id"`
	KratosIdentityID          string    `json:"kratos_identity_id"`
	Email                     string    `json:"email"`
	FullName                  *string   `json:"full_name,omitempty"`
	CPF                       *string   `json:"cpf,omitempty"`
	BalanceCents              int64     `json:"balance_cents"`
	DailyTransferLimitCents   int64     `json:"daily_transfer_limit_cents"`
	MonthlyTransferLimitCents int64     `json:"monthly_transfer_limit_cents"`
	Status                    string    `json:"status"`
	KYCStatus                 string    `json:"kyc_status"`
	CreatedAt                 time.Time `json:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at"`
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	KratosIdentityID string  `json:"kratos_identity_id"`
	Email            string  `json:"email"`
	FullName         *string `json:"full_name,omitempty"`
	CPF              *string `json:"cpf,omitempty"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	FullName *string `json:"full_name,omitempty"`
}
