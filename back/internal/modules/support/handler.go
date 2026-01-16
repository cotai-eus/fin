package support

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/lauratech/fin/back/internal/shared/response"

	"github.com/go-chi/chi/v5"
)

// Handler handles HTTP requests for support tickets
type Handler struct {
	service *Service
}

// NewHandler creates a new support handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateTicket handles POST /api/support/tickets
func (h *Handler) CreateTicket(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Decode request body
	var req CreateTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	// Create ticket
	ticket, err := h.service.CreateTicket(r.Context(), userID, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Return success response
	response.Success(w, http.StatusCreated, ticket, r.Context())
}

// ListTickets handles GET /api/support/tickets
func (h *Handler) ListTickets(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Parse pagination parameters
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	// Get optional status filter
	status := r.URL.Query().Get("status")

	// List tickets
	params := ListTicketsParams{
		Page:   page,
		Limit:  limit,
		Status: status,
	}

	tickets, total, err := h.service.ListUserTickets(r.Context(), userID, params)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Calculate pagination metadata
	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	pagination := response.Pagination{
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
		HasMore:    page < totalPages,
	}

	// Return paginated response
	response.Paginated(w, http.StatusOK, tickets, pagination, r.Context())
}

// GetTicket handles GET /api/support/tickets/{id}
func (h *Handler) GetTicket(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Get ticket ID from URL
	ticketID := chi.URLParam(r, "id")
	if ticketID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Ticket ID is required", nil)
		return
	}

	// Get ticket
	ticket, err := h.service.GetTicketByID(r.Context(), userID, ticketID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Return success response
	response.Success(w, http.StatusOK, ticket, r.Context())
}

// GetTicketWithMessages handles GET /api/support/tickets/{id}/messages
func (h *Handler) GetTicketWithMessages(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Get ticket ID from URL
	ticketID := chi.URLParam(r, "id")
	if ticketID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Ticket ID is required", nil)
		return
	}

	// Parse pagination parameters for messages
	page := 1
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 50 // Higher limit for messages
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	params := ListMessagesParams{
		Page:  page,
		Limit: limit,
	}

	// Get ticket with messages
	ticketWithMessages, err := h.service.GetTicketWithMessages(r.Context(), userID, ticketID, params)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Return success response
	response.Success(w, http.StatusOK, ticketWithMessages, r.Context())
}

// AddMessage handles POST /api/support/tickets/{id}/messages
func (h *Handler) AddMessage(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Get ticket ID from URL
	ticketID := chi.URLParam(r, "id")
	if ticketID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Ticket ID is required", nil)
		return
	}

	// Decode request body
	var req AddMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	// Add message
	message, err := h.service.AddMessage(r.Context(), userID, ticketID, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Return success response
	response.Success(w, http.StatusCreated, message, r.Context())
}

// UpdateTicketStatus handles PATCH /api/support/tickets/{id}/status
func (h *Handler) UpdateTicketStatus(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Get ticket ID from URL
	ticketID := chi.URLParam(r, "id")
	if ticketID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Ticket ID is required", nil)
		return
	}

	// Decode request body
	var req UpdateTicketStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	// Update status
	ticket, err := h.service.UpdateTicketStatus(r.Context(), userID, ticketID, req)
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Return success response
	response.Success(w, http.StatusOK, ticket, r.Context())
}

// GetTicketStats handles GET /api/support/stats (admin only - for future implementation)
func (h *Handler) GetTicketStats(w http.ResponseWriter, r *http.Request) {
	// Get ticket statistics
	stats, err := h.service.GetTicketStats(r.Context())
	if err != nil {
		h.handleError(w, err)
		return
	}

	// Return success response
	response.Success(w, http.StatusOK, stats, r.Context())
}

// handleError maps domain errors to HTTP responses
func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch err {
	case ErrTicketNotFound:
		response.Error(w, http.StatusNotFound, "TICKET_001", "Ticket not found", nil)
	case ErrMessageNotFound:
		response.Error(w, http.StatusNotFound, "TICKET_002", "Message not found", nil)
	case ErrUnauthorizedAccess:
		response.Error(w, http.StatusForbidden, "AUTH_002", "Unauthorized access to ticket", nil)
	case ErrInvalidCategory:
		response.Error(w, http.StatusBadRequest, "VAL_101", "Invalid ticket category", nil)
	case ErrInvalidPriority:
		response.Error(w, http.StatusBadRequest, "VAL_102", "Invalid ticket priority", nil)
	case ErrInvalidStatus:
		response.Error(w, http.StatusBadRequest, "VAL_103", "Invalid ticket status", nil)
	case ErrInvalidStatusTransition:
		response.Error(w, http.StatusBadRequest, "BUS_101", "Invalid status transition", nil)
	case ErrEmptySubject:
		response.Error(w, http.StatusBadRequest, "VAL_104", "Ticket subject cannot be empty", nil)
	case ErrEmptyDescription:
		response.Error(w, http.StatusBadRequest, "VAL_105", "Ticket description cannot be empty", nil)
	case ErrEmptyMessage:
		response.Error(w, http.StatusBadRequest, "VAL_106", "Message cannot be empty", nil)
	case ErrTicketClosed:
		response.Error(w, http.StatusBadRequest, "BUS_102", "Cannot add message to closed ticket", nil)
	case ErrSubjectTooLong:
		response.Error(w, http.StatusBadRequest, "VAL_107", "Subject exceeds maximum length", nil)
	case ErrDescriptionTooShort:
		response.Error(w, http.StatusBadRequest, "VAL_108", "Description too short (minimum 10 characters)", nil)
	case ErrMessageTooShort:
		response.Error(w, http.StatusBadRequest, "VAL_109", "Message too short (minimum 5 characters)", nil)
	default:
		response.Error(w, http.StatusInternalServerError, "SYS_001", "Internal server error", nil)
	}
}
