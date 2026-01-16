package support

import (
	"context"
	"database/sql"
	"errors"

	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"

	"github.com/google/uuid"
)

// Repository handles data access for support tickets
type Repository struct {
	db      *sql.DB
	queries *db.Queries
}

// NewRepository creates a new support repository
func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db:      database,
		queries: db.New(database),
	}
}

// CreateTicket creates a new support ticket
func (r *Repository) CreateTicket(ctx context.Context, params db.CreateTicketParams) (*db.SupportTicket, error) {
	ticket, err := r.queries.CreateTicket(ctx, params)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// GetByID retrieves a ticket by ID
func (r *Repository) GetByID(ctx context.Context, id string) (*db.SupportTicket, error) {
	ticketID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	ticket, err := r.queries.GetTicketByID(ctx, ticketID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	return &ticket, nil
}

// GetByNumber retrieves a ticket by ticket number
func (r *Repository) GetByNumber(ctx context.Context, ticketNumber string) (*db.SupportTicket, error) {
	ticket, err := r.queries.GetTicketByNumber(ctx, ticketNumber)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	return &ticket, nil
}

// ListUserTickets retrieves all tickets for a user with pagination
func (r *Repository) ListUserTickets(ctx context.Context, userID string, limit, offset int32) ([]db.SupportTicket, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	tickets, err := r.queries.ListUserTickets(ctx, db.ListUserTicketsParams{
		UserID: userUUID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

// CountUserTickets counts total tickets for a user
func (r *Repository) CountUserTickets(ctx context.Context, userID string) (int64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}

	count, err := r.queries.CountUserTickets(ctx, userUUID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ListUserTicketsByStatus retrieves user tickets filtered by status
func (r *Repository) ListUserTicketsByStatus(ctx context.Context, userID, status string, limit, offset int32) ([]db.SupportTicket, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	tickets, err := r.queries.ListUserTicketsByStatus(ctx, db.ListUserTicketsByStatusParams{
		UserID: userUUID,
		Status: status,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

// CountUserTicketsByStatus counts user tickets by status
func (r *Repository) CountUserTicketsByStatus(ctx context.Context, userID, status string) (int64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}

	count, err := r.queries.CountUserTicketsByStatus(ctx, db.CountUserTicketsByStatusParams{
		UserID: userUUID,
		Status: status,
	})
	if err != nil {
		return 0, err
	}

	return count, nil
}

// UpdateStatus updates the status of a ticket
func (r *Repository) UpdateStatus(ctx context.Context, id, status string) (*db.SupportTicket, error) {
	ticketID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	ticket, err := r.queries.UpdateTicketStatus(ctx, db.UpdateTicketStatusParams{
		ID:     ticketID,
		Status: status,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	return &ticket, nil
}

// UpdateTicket updates ticket status and priority
func (r *Repository) UpdateTicket(ctx context.Context, id, status, priority string) (*db.SupportTicket, error) {
	ticketID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	ticket, err := r.queries.UpdateTicket(ctx, db.UpdateTicketParams{
		ID:       ticketID,
		Status:   status,
		Priority: priority,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	return &ticket, nil
}

// GetForUpdate retrieves a ticket with FOR UPDATE lock
func (r *Repository) GetForUpdate(ctx context.Context, id string) (*db.SupportTicket, error) {
	ticketID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	ticket, err := r.queries.GetTicketForUpdate(ctx, ticketID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	return &ticket, nil
}

// Delete deletes a ticket
func (r *Repository) Delete(ctx context.Context, id string) error {
	ticketID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteTicket(ctx, ticketID)
}

// CreateMessage creates a new ticket message
func (r *Repository) CreateMessage(ctx context.Context, params db.CreateTicketMessageParams) (*db.TicketMessage, error) {
	message, err := r.queries.CreateTicketMessage(ctx, params)
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// GetMessageByID retrieves a message by ID
func (r *Repository) GetMessageByID(ctx context.Context, id string) (*db.TicketMessage, error) {
	messageID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	message, err := r.queries.GetTicketMessageByID(ctx, messageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}

	return &message, nil
}

// ListMessages retrieves messages for a ticket with pagination
func (r *Repository) ListMessages(ctx context.Context, ticketID string, limit, offset int32) ([]db.TicketMessage, error) {
	ticketUUID, err := uuid.Parse(ticketID)
	if err != nil {
		return nil, err
	}

	messages, err := r.queries.ListTicketMessages(ctx, db.ListTicketMessagesParams{
		TicketID: ticketUUID,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, err
	}

	return messages, nil
}

// CountMessages counts messages for a ticket
func (r *Repository) CountMessages(ctx context.Context, ticketID string) (int64, error) {
	ticketUUID, err := uuid.Parse(ticketID)
	if err != nil {
		return 0, err
	}

	count, err := r.queries.CountTicketMessages(ctx, ticketUUID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetLatestMessage retrieves the latest message for a ticket
func (r *Repository) GetLatestMessage(ctx context.Context, ticketID string) (*db.TicketMessage, error) {
	ticketUUID, err := uuid.Parse(ticketID)
	if err != nil {
		return nil, err
	}

	message, err := r.queries.GetLatestTicketMessage(ctx, ticketUUID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrMessageNotFound
		}
		return nil, err
	}

	return &message, nil
}

// DeleteMessage deletes a message
func (r *Repository) DeleteMessage(ctx context.Context, id string) error {
	messageID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteTicketMessage(ctx, messageID)
}

// Admin/Staff Methods

// ListAllTickets retrieves all tickets (admin view)
func (r *Repository) ListAllTickets(ctx context.Context, limit, offset int32) ([]db.SupportTicket, error) {
	tickets, err := r.queries.ListAllTickets(ctx, db.ListAllTicketsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

// CountAllTickets counts all tickets
func (r *Repository) CountAllTickets(ctx context.Context) (int64, error) {
	count, err := r.queries.CountAllTickets(ctx)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ListTicketsByStatus retrieves tickets filtered by status (admin view)
func (r *Repository) ListTicketsByStatus(ctx context.Context, status string, limit, offset int32) ([]db.SupportTicket, error) {
	tickets, err := r.queries.ListTicketsByStatus(ctx, db.ListTicketsByStatusParams{
		Status: status,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

// CountTicketsByStatus counts tickets by status
func (r *Repository) CountTicketsByStatus(ctx context.Context, status string) (int64, error) {
	count, err := r.queries.CountTicketsByStatus(ctx, status)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetTicketStats retrieves ticket statistics
func (r *Repository) GetTicketStats(ctx context.Context) (*db.GetTicketStatsRow, error) {
	stats, err := r.queries.GetTicketStats(ctx)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
