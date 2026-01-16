package support

import (
	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

// dbTicketToTicket converts a database ticket to domain ticket
func dbTicketToTicket(dbTicket *db.SupportTicket) *Ticket {
	if dbTicket == nil {
		return nil
	}

	return &Ticket{
		ID:           dbTicket.ID.String(),
		UserID:       dbTicket.UserID.String(),
		TicketNumber: dbTicket.TicketNumber,
		Category:     dbTicket.Category,
		Priority:     dbTicket.Priority,
		Status:       dbTicket.Status,
		Subject:      dbTicket.Subject,
		Description:  dbTicket.Description,
		CreatedAt:    dbTicket.CreatedAt.Time,
		UpdatedAt:    dbTicket.UpdatedAt.Time,
	}
}

// dbTicketToTicketSummary converts a database ticket to ticket summary
func dbTicketToTicketSummary(dbTicket *db.SupportTicket) *TicketSummary {
	if dbTicket == nil {
		return nil
	}

	return &TicketSummary{
		ID:           dbTicket.ID.String(),
		TicketNumber: dbTicket.TicketNumber,
		Category:     dbTicket.Category,
		Priority:     dbTicket.Priority,
		Status:       dbTicket.Status,
		Subject:      dbTicket.Subject,
		CreatedAt:    dbTicket.CreatedAt.Time,
		UpdatedAt:    dbTicket.UpdatedAt.Time,
	}
}

// dbTicketsToTicketSummaries converts a slice of database tickets to ticket summaries
func dbTicketsToTicketSummaries(dbTickets []db.SupportTicket) []*TicketSummary {
	summaries := make([]*TicketSummary, len(dbTickets))
	for i, dbTicket := range dbTickets {
		summaries[i] = dbTicketToTicketSummary(&dbTicket)
	}
	return summaries
}

// dbMessageToTicketMessage converts a database ticket message to domain message
func dbMessageToTicketMessage(dbMessage *db.TicketMessage) *TicketMessage {
	if dbMessage == nil {
		return nil
	}

	return &TicketMessage{
		ID:        dbMessage.ID.String(),
		TicketID:  dbMessage.TicketID.String(),
		UserID:    dbMessage.UserID.String(),
		Message:   dbMessage.Message,
		IsStaff:   dbMessage.IsStaff.Bool,
		CreatedAt: dbMessage.CreatedAt.Time,
	}
}

// dbMessagesToTicketMessages converts a slice of database messages to domain messages
func dbMessagesToTicketMessages(dbMessages []db.TicketMessage) []*TicketMessage {
	messages := make([]*TicketMessage, len(dbMessages))
	for i, dbMessage := range dbMessages {
		messages[i] = dbMessageToTicketMessage(&dbMessage)
	}
	return messages
}
