package support

import "time"

// Ticket represents a support ticket domain model
type Ticket struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	TicketNumber string    `json:"ticket_number"`
	Category     string    `json:"category"`
	Priority     string    `json:"priority"`
	Status       string    `json:"status"`
	Subject      string    `json:"subject"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TicketMessage represents a message in a support ticket
type TicketMessage struct {
	ID        string    `json:"id"`
	TicketID  string    `json:"ticket_id"`
	UserID    string    `json:"user_id"`
	Message   string    `json:"message"`
	IsStaff   bool      `json:"is_staff"`
	CreatedAt time.Time `json:"created_at"`
}

// TicketWithMessages represents a ticket with its message thread
type TicketWithMessages struct {
	Ticket   *Ticket          `json:"ticket"`
	Messages []*TicketMessage `json:"messages"`
}

// TicketSummary represents a lightweight ticket for list views
type TicketSummary struct {
	ID           string    `json:"id"`
	TicketNumber string    `json:"ticket_number"`
	Category     string    `json:"category"`
	Priority     string    `json:"priority"`
	Status       string    `json:"status"`
	Subject      string    `json:"subject"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TicketStats represents ticket statistics for admin dashboard
type TicketStats struct {
	OpenCount       int64 `json:"open_count"`
	InProgressCount int64 `json:"in_progress_count"`
	WaitingCount    int64 `json:"waiting_count"`
	ResolvedCount   int64 `json:"resolved_count"`
	ClosedCount     int64 `json:"closed_count"`
	UrgentCount     int64 `json:"urgent_count"`
	HighCount       int64 `json:"high_count"`
	TotalCount      int64 `json:"total_count"`
}

// CreateTicketRequest represents the request to create a new ticket
type CreateTicketRequest struct {
	Category    string `json:"category"`
	Priority    string `json:"priority"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

// UpdateTicketStatusRequest represents the request to update ticket status
type UpdateTicketStatusRequest struct {
	Status string `json:"status"`
}

// UpdateTicketRequest represents the request to update ticket details
type UpdateTicketRequest struct {
	Status   string `json:"status"`
	Priority string `json:"priority"`
}

// AddMessageRequest represents the request to add a message to a ticket
type AddMessageRequest struct {
	Message string `json:"message"`
}

// ListTicketsParams represents parameters for listing tickets
type ListTicketsParams struct {
	Page   int
	Limit  int
	Status string
}

// ListMessagesParams represents parameters for listing messages
type ListMessagesParams struct {
	Page  int
	Limit int
}

// Valid ticket categories
var ValidCategories = map[string]bool{
	"account":   true,
	"card":      true,
	"transfer":  true,
	"bill":      true,
	"technical": true,
	"other":     true,
}

// Valid ticket priorities
var ValidPriorities = map[string]bool{
	"low":    true,
	"medium": true,
	"high":   true,
	"urgent": true,
}

// Valid ticket statuses
var ValidStatuses = map[string]bool{
	"open":        true,
	"in_progress": true,
	"waiting":     true,
	"resolved":    true,
	"closed":      true,
}

// AllowedStatusTransitions defines allowed status transitions
var AllowedStatusTransitions = map[string][]string{
	"open":        {"in_progress", "waiting", "resolved", "closed"},
	"in_progress": {"waiting", "resolved", "closed", "open"},
	"waiting":     {"in_progress", "resolved", "closed"},
	"resolved":    {"closed", "open"}, // Can reopen if needed
	"closed":      {"open"},           // Admin can reopen
}
