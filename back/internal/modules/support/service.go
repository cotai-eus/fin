package support

import (
	"context"
	"database/sql"
	"fmt"
	"sync/atomic"
	"time"

	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"

	"github.com/google/uuid"
)

// Service handles business logic for support tickets
type Service struct {
	repo    *Repository
	db      *sql.DB
	counter atomic.Int64
}

// NewService creates a new support service
func NewService(repo *Repository, database *sql.DB) *Service {
	return &Service{
		repo: repo,
		db:   database,
	}
}

// generateTicketNumber generates a unique ticket number
// Format: SUP-YYYYMMDD-####
func (s *Service) generateTicketNumber() string {
	date := time.Now().Format("20060102")
	count := s.counter.Add(1)
	return fmt.Sprintf("SUP-%s-%04d", date, count%10000)
}

// CreateTicket creates a new support ticket
func (s *Service) CreateTicket(ctx context.Context, userID string, req CreateTicketRequest) (*Ticket, error) {
	// Validate request
	if err := ValidateCreateTicketRequest(req); err != nil {
		return nil, err
	}

	// Parse user ID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Generate ticket number
	ticketNumber := s.generateTicketNumber()

	// Create ticket
	params := db.CreateTicketParams{
		UserID:       userUUID,
		TicketNumber: ticketNumber,
		Category:     req.Category,
		Priority:     req.Priority,
		Status:       "open",
		Subject:      req.Subject,
		Description:  req.Description,
	}

	dbTicket, err := s.repo.CreateTicket(ctx, params)
	if err != nil {
		return nil, err
	}

	return dbTicketToTicket(dbTicket), nil
}

// GetTicketByID retrieves a ticket by ID with ownership verification
func (s *Service) GetTicketByID(ctx context.Context, userID, ticketID string) (*Ticket, error) {
	dbTicket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if dbTicket.UserID.String() != userID {
		return nil, ErrUnauthorizedAccess
	}

	return dbTicketToTicket(dbTicket), nil
}

// GetTicketByNumber retrieves a ticket by ticket number with ownership verification
func (s *Service) GetTicketByNumber(ctx context.Context, userID, ticketNumber string) (*Ticket, error) {
	dbTicket, err := s.repo.GetByNumber(ctx, ticketNumber)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if dbTicket.UserID.String() != userID {
		return nil, ErrUnauthorizedAccess
	}

	return dbTicketToTicket(dbTicket), nil
}

// ListUserTickets retrieves all tickets for a user with pagination
func (s *Service) ListUserTickets(ctx context.Context, userID string, params ListTicketsParams) ([]*TicketSummary, int64, error) {
	limit := int32(params.Limit)
	offset := int32((params.Page - 1) * params.Limit)

	var dbTickets []db.SupportTicket
	var total int64
	var err error

	// Filter by status if provided
	if params.Status != "" {
		if err := ValidateStatus(params.Status); err != nil {
			return nil, 0, err
		}

		dbTickets, err = s.repo.ListUserTicketsByStatus(ctx, userID, params.Status, limit, offset)
		if err != nil {
			return nil, 0, err
		}

		total, err = s.repo.CountUserTicketsByStatus(ctx, userID, params.Status)
		if err != nil {
			return nil, 0, err
		}
	} else {
		dbTickets, err = s.repo.ListUserTickets(ctx, userID, limit, offset)
		if err != nil {
			return nil, 0, err
		}

		total, err = s.repo.CountUserTickets(ctx, userID)
		if err != nil {
			return nil, 0, err
		}
	}

	summaries := dbTicketsToTicketSummaries(dbTickets)
	return summaries, total, nil
}

// UpdateTicketStatus updates the status of a ticket with ownership verification
func (s *Service) UpdateTicketStatus(ctx context.Context, userID, ticketID string, req UpdateTicketStatusRequest) (*Ticket, error) {
	// Validate new status
	if err := ValidateStatus(req.Status); err != nil {
		return nil, err
	}

	// Execute in transaction to ensure atomicity
	var updatedTicket *db.SupportTicket
	err := s.executeInTransaction(ctx, func(tx *sql.Tx) error {
		qtx := db.New(tx)

		// Get ticket with lock
		ticketUUID, err := uuid.Parse(ticketID)
		if err != nil {
			return err
		}

		dbTicket, err := qtx.GetTicketForUpdate(ctx, ticketUUID)
		if err != nil {
			return err
		}

		// Verify ownership
		if dbTicket.UserID.String() != userID {
			return ErrUnauthorizedAccess
		}

		// Validate status transition
		if err := ValidateStatusTransition(dbTicket.Status, req.Status); err != nil {
			return err
		}

		// Update status
		updated, err := qtx.UpdateTicketStatus(ctx, db.UpdateTicketStatusParams{
			ID:     ticketUUID,
			Status: req.Status,
		})
		if err != nil {
			return err
		}

		updatedTicket = &updated
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dbTicketToTicket(updatedTicket), nil
}

// AddMessage adds a message to a ticket with ownership verification
func (s *Service) AddMessage(ctx context.Context, userID, ticketID string, req AddMessageRequest) (*TicketMessage, error) {
	// Validate message
	if err := ValidateMessage(req.Message); err != nil {
		return nil, err
	}

	// Get ticket and verify ownership
	dbTicket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	if dbTicket.UserID.String() != userID {
		return nil, ErrUnauthorizedAccess
	}

	// Check if ticket is closed
	if dbTicket.Status == "closed" {
		return nil, ErrTicketClosed
	}

	// Parse IDs
	ticketUUID, err := uuid.Parse(ticketID)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Create message
	params := db.CreateTicketMessageParams{
		TicketID: ticketUUID,
		UserID:   userUUID,
		Message:  req.Message,
		IsStaff:  sql.NullBool{Bool: false, Valid: true},
	}

	dbMessage, err := s.repo.CreateMessage(ctx, params)
	if err != nil {
		return nil, err
	}

	return dbMessageToTicketMessage(dbMessage), nil
}

// GetTicketWithMessages retrieves a ticket with all its messages
func (s *Service) GetTicketWithMessages(ctx context.Context, userID, ticketID string, messagesParams ListMessagesParams) (*TicketWithMessages, error) {
	// Get ticket
	ticket, err := s.GetTicketByID(ctx, userID, ticketID)
	if err != nil {
		return nil, err
	}

	// Get messages
	limit := int32(messagesParams.Limit)
	offset := int32((messagesParams.Page - 1) * messagesParams.Limit)

	dbMessages, err := s.repo.ListMessages(ctx, ticketID, limit, offset)
	if err != nil {
		return nil, err
	}

	messages := dbMessagesToTicketMessages(dbMessages)

	return &TicketWithMessages{
		Ticket:   ticket,
		Messages: messages,
	}, nil
}

// ListMessages retrieves messages for a ticket with pagination
func (s *Service) ListMessages(ctx context.Context, userID, ticketID string, params ListMessagesParams) ([]*TicketMessage, int64, error) {
	// Verify ticket ownership
	dbTicket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return nil, 0, err
	}

	if dbTicket.UserID.String() != userID {
		return nil, 0, ErrUnauthorizedAccess
	}

	// Get messages
	limit := int32(params.Limit)
	offset := int32((params.Page - 1) * params.Limit)

	dbMessages, err := s.repo.ListMessages(ctx, ticketID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.CountMessages(ctx, ticketID)
	if err != nil {
		return nil, 0, err
	}

	messages := dbMessagesToTicketMessages(dbMessages)
	return messages, total, nil
}

// DeleteTicket soft deletes a ticket (for future implementation)
// Currently not implemented as per requirements
func (s *Service) DeleteTicket(ctx context.Context, userID, ticketID string) error {
	// Get ticket and verify ownership
	dbTicket, err := s.repo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}

	if dbTicket.UserID.String() != userID {
		return ErrUnauthorizedAccess
	}

	// Delete ticket
	return s.repo.Delete(ctx, ticketID)
}

// executeInTransaction executes a function within a database transaction
func (s *Service) executeInTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Admin/Staff Methods (for future implementation)

// GetTicketStats retrieves ticket statistics (admin only)
func (s *Service) GetTicketStats(ctx context.Context) (*TicketStats, error) {
	dbStats, err := s.repo.GetTicketStats(ctx)
	if err != nil {
		return nil, err
	}

	stats := &TicketStats{
		OpenCount:       dbStats.OpenCount,
		InProgressCount: dbStats.InProgressCount,
		WaitingCount:    dbStats.WaitingCount,
		ResolvedCount:   dbStats.ResolvedCount,
		ClosedCount:     dbStats.ClosedCount,
		UrgentCount:     dbStats.UrgentCount,
		HighCount:       dbStats.HighCount,
		TotalCount:      dbStats.TotalCount,
	}

	return stats, nil
}
