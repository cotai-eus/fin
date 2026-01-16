package transfers

import "time"

// Transfer represents a transfer in the system
type Transfer struct {
	ID              string     `json:"id"`
	UserID          string     `json:"user_id"`
	Type            string     `json:"type"`   // "pix", "ted", "p2p", "deposit", "withdrawal"
	Status          string     `json:"status"` // "pending", "processing", "completed", "failed", "cancelled"
	AmountCents     int64      `json:"amount_cents"`
	FeeCents        int64      `json:"fee_cents"`
	Currency        string     `json:"currency"`
	PixKey          *string    `json:"pix_key,omitempty"`
	PixKeyType      *string    `json:"pix_key_type,omitempty"`
	RecipientName   *string    `json:"recipient_name,omitempty"`
	RecipientDocument *string  `json:"recipient_document,omitempty"`
	RecipientBank   *string    `json:"recipient_bank,omitempty"`
	RecipientBranch *string    `json:"recipient_branch,omitempty"`
	RecipientAccount *string   `json:"recipient_account,omitempty"`
	RecipientAccountType *string `json:"recipient_account_type,omitempty"`
	RecipientUserID *string    `json:"recipient_user_id,omitempty"`
	ScheduledFor    *time.Time `json:"scheduled_for,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	FailureReason   *string    `json:"failure_reason,omitempty"`
	AuthenticationCode *string `json:"authentication_code,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CreatePIXRequest represents a request to create a PIX transfer
type CreatePIXRequest struct {
	PixKey      string `json:"pix_key"`
	PixKeyType  string `json:"pix_key_type"` // "cpf", "cnpj", "email", "phone", "random"
	AmountCents int64  `json:"amount_cents"`
	Description string `json:"description,omitempty"`
}

// CreateTEDRequest represents a request to create a TED transfer
type CreateTEDRequest struct {
	RecipientName        string `json:"recipient_name"`
	RecipientDocument    string `json:"recipient_document"`
	RecipientBank        string `json:"recipient_bank"`    // 3 digits
	RecipientBranch      string `json:"recipient_branch"`  // 4-5 digits
	RecipientAccount     string `json:"recipient_account"` // up to 12 digits
	RecipientAccountType string `json:"recipient_account_type"` // "checking", "savings"
	AmountCents          int64  `json:"amount_cents"`
	Description          string `json:"description,omitempty"`
}

// CreateP2PRequest represents a request to create a P2P (peer-to-peer) transfer
type CreateP2PRequest struct {
	RecipientUserID string `json:"recipient_user_id"`
	AmountCents     int64  `json:"amount_cents"`
	Description     string `json:"description,omitempty"`
}

// TransferListParams represents pagination parameters for listing transfers
type TransferListParams struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}
