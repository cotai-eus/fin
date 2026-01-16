package bills

import "time"

// Bill represents a bill in the system
type Bill struct {
	ID               string     `json:"id"`
	UserID           string     `json:"user_id"`
	Type             string     `json:"type"`   // "bank", "utility", "tax", "other"
	Status           string     `json:"status"` // "pending", "paid", "overdue", "cancelled"
	Barcode          string     `json:"barcode"`
	AmountCents      int64      `json:"amount_cents"`
	FeeCents         int64      `json:"fee_cents"`
	FinalAmountCents int64      `json:"final_amount_cents"`
	RecipientName    string     `json:"recipient_name"`
	DueDate          time.Time  `json:"due_date"`
	PaymentDate      *time.Time `json:"payment_date,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
}

// BillSummary represents a bill summary for list responses
type BillSummary struct {
	ID               string    `json:"id"`
	Type             string    `json:"type"`
	Status           string    `json:"status"`
	RecipientName    string    `json:"recipient_name"`
	AmountCents      int64     `json:"amount_cents"`
	FinalAmountCents int64     `json:"final_amount_cents"`
	DueDate          time.Time `json:"due_date"`
	CreatedAt        time.Time `json:"created_at"`
}

// ValidateBarcodeRequest represents a barcode validation request
type ValidateBarcodeRequest struct {
	Barcode string `json:"barcode"`
}

// ValidateBarcodeResponse represents a barcode validation response
type ValidateBarcodeResponse struct {
	Valid         bool   `json:"valid"`
	RecipientName string `json:"recipient_name,omitempty"`
	AmountCents   int64  `json:"amount_cents,omitempty"`
	DueDate       string `json:"due_date,omitempty"`
	Type          string `json:"type,omitempty"`
}

// CreateBillRequest represents a request to create/register a bill
type CreateBillRequest struct {
	Barcode string `json:"barcode"`
	Type    string `json:"type"` // "bank", "utility", "tax", "other"
}

// PayBillRequest represents a request to pay a bill
type PayBillRequest struct {
	BillID string `json:"bill_id"`
}

// BillListParams represents pagination parameters for listing bills
type BillListParams struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Status string `json:"status,omitempty"` // Filter by status
}
