package cards

import "time"

// Card represents a card with sensitive data (internal use only)
// This struct contains decrypted card data and should NEVER be serialized to JSON directly
type Card struct {
	ID                       string     `json:"id"`
	UserID                   string     `json:"user_id"`
	Type                     string     `json:"type"`
	Brand                    string     `json:"brand"`
	Status                   string     `json:"status"`
	CardNumber               string     `json:"-"` // NEVER expose in JSON - security critical
	CVV                      string     `json:"-"` // NEVER expose in JSON - security critical
	LastFourDigits           string     `json:"last_four_digits"`
	HolderName               string     `json:"holder_name"`
	ExpiryMonth              int        `json:"expiry_month"`
	ExpiryYear               int        `json:"expiry_year"`
	DailyLimitCents          int64      `json:"daily_limit_cents"`
	MonthlyLimitCents        int64      `json:"monthly_limit_cents"`
	CurrentDailySpentCents   int64      `json:"current_daily_spent_cents"`
	CurrentMonthlySpentCents int64      `json:"current_monthly_spent_cents"`
	IsContactless            bool       `json:"is_contactless"`
	IsInternational          bool       `json:"is_international"`
	BlockInternational       bool       `json:"block_international"`
	BlockOnline              bool       `json:"block_online"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	ExpiresAt                time.Time  `json:"expires_at"`
	BlockedAt                *time.Time `json:"blocked_at,omitempty"`
}

// CardSummary represents card info for listing (no sensitive data)
// Use this for GET /cards list responses
type CardSummary struct {
	ID                       string    `json:"id"`
	UserID                   string    `json:"user_id"`
	Type                     string    `json:"type"`
	Brand                    string    `json:"brand"`
	Status                   string    `json:"status"`
	LastFourDigits           string    `json:"last_four_digits"`
	HolderName               string    `json:"holder_name"`
	ExpiryMonth              int       `json:"expiry_month"`
	ExpiryYear               int       `json:"expiry_year"`
	DailyLimitCents          int64     `json:"daily_limit_cents"`
	MonthlyLimitCents        int64     `json:"monthly_limit_cents"`
	CurrentDailySpentCents   int64     `json:"current_daily_spent_cents"`
	CurrentMonthlySpentCents int64     `json:"current_monthly_spent_cents"`
	IsContactless            bool      `json:"is_contactless"`
	IsInternational          bool      `json:"is_international"`
	CreatedAt                time.Time `json:"created_at"`
}

// CardDetails represents card info with masked number (for GET /cards/{id})
// Use this for individual card detail responses
type CardDetails struct {
	ID                       string     `json:"id"`
	Type                     string     `json:"type"`
	Brand                    string     `json:"brand"`
	Status                   string     `json:"status"`
	LastFourDigits           string     `json:"last_four_digits"`
	MaskedCardNumber         string     `json:"masked_card_number"` // e.g., "**** **** **** 1234"
	HolderName               string     `json:"holder_name"`
	ExpiryMonth              int        `json:"expiry_month"`
	ExpiryYear               int        `json:"expiry_year"`
	DailyLimitCents          int64      `json:"daily_limit_cents"`
	MonthlyLimitCents        int64      `json:"monthly_limit_cents"`
	CurrentDailySpentCents   int64      `json:"current_daily_spent_cents"`
	CurrentMonthlySpentCents int64      `json:"current_monthly_spent_cents"`
	IsContactless            bool       `json:"is_contactless"`
	IsInternational          bool       `json:"is_international"`
	BlockInternational       bool       `json:"block_international"`
	BlockOnline              bool       `json:"block_online"`
	HasPIN                   bool       `json:"has_pin"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	ExpiresAt                time.Time  `json:"expires_at"`
	BlockedAt                *time.Time `json:"blocked_at,omitempty"`
}

// CreateCardRequest for POST /api/cards
type CreateCardRequest struct {
	Type              string `json:"type" validate:"required,oneof=physical virtual"`
	Brand             string `json:"brand" validate:"required,oneof=visa mastercard elo"`
	CardNumber        string `json:"card_number,omitempty"` // Optional - auto-generate if empty
	CVV               string `json:"cvv" validate:"required"`
	PIN               string `json:"pin,omitempty"` // Optional on creation
	HolderName        string `json:"holder_name" validate:"required"`
	DailyLimitCents   int64  `json:"daily_limit_cents,omitempty"`
	MonthlyLimitCents int64  `json:"monthly_limit_cents,omitempty"`
}

// UpdateLimitsRequest for PATCH /api/cards/{id}/limits
type UpdateLimitsRequest struct {
	DailyLimitCents   int64 `json:"daily_limit_cents" validate:"required,min=0"`
	MonthlyLimitCents int64 `json:"monthly_limit_cents" validate:"required,min=0"`
}

// SecuritySettingsRequest for PATCH /api/cards/{id}/security
type SecuritySettingsRequest struct {
	IsContactless      bool `json:"is_contactless"`
	IsInternational    bool `json:"is_international"`
	BlockInternational bool `json:"block_international"`
	BlockOnline        bool `json:"block_online"`
}

// SetPINRequest for POST /api/cards/{id}/pin
type SetPINRequest struct {
	PIN        string `json:"pin" validate:"required"`
	CurrentPIN string `json:"current_pin,omitempty"` // Required if changing existing PIN
}

// CancelCardRequest for DELETE /api/cards/{id}
type CancelCardRequest struct {
	Reason string `json:"reason" validate:"required,oneof=lost stolen damaged user_request"`
}

// CreateCardParams holds parameters for creating a card in the repository
type CreateCardParams struct {
	UserID             string
	Type               string
	Brand              string
	CardNumber         string
	CVV                string
	PIN                string
	HolderName         string
	ExpiryMonth        int
	ExpiryYear         int
	DailyLimitCents    int64
	MonthlyLimitCents  int64
	IsContactless      bool
	IsInternational    bool
	BlockInternational bool
	BlockOnline        bool
	ExpiresAt          time.Time
}

// UpdateLimitsParams holds parameters for updating card limits
type UpdateLimitsParams struct {
	DailyLimitCents   int64
	MonthlyLimitCents int64
}

// SecuritySettingsParams holds parameters for updating security settings
type SecuritySettingsParams struct {
	IsContactless      bool
	IsInternational    bool
	BlockInternational bool
	BlockOnline        bool
}
