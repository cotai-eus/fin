package transfers

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lauratech/fin/back/internal/shared/response"
)

// Handler handles HTTP requests for transfers
type Handler struct {
	service *Service
}

// NewHandler creates a new transfer handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// ExecutePIX executes a PIX transfer
// POST /api/transfers/pix
func (h *Handler) ExecutePIX(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Decode request
	var req CreatePIXRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	// Execute transfer
	transfer, err := h.service.ExecutePIX(r.Context(), userID, req)
	if err != nil {
		h.handleTransferError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, transfer, r.Context())
}

// ExecuteTED executes a TED transfer
// POST /api/transfers/ted
func (h *Handler) ExecuteTED(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Decode request
	var req CreateTEDRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	// Execute transfer
	transfer, err := h.service.ExecuteTED(r.Context(), userID, req)
	if err != nil {
		h.handleTransferError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, transfer, r.Context())
}

// ExecuteP2P executes a peer-to-peer transfer
// POST /api/transfers/p2p
func (h *Handler) ExecuteP2P(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Decode request
	var req CreateP2PRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	// Execute transfer
	transfer, err := h.service.ExecuteP2P(r.Context(), userID, req)
	if err != nil {
		h.handleTransferError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, transfer, r.Context())
}

// List retrieves transfers for the current user with pagination
// GET /api/transfers?page=1&limit=20
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Parse query parameters
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

	// Get transfers
	transfers, total, err := h.service.List(r.Context(), userID, TransferListParams{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "SYS_001", "Internal server error", nil)
		return
	}

	// Calculate pagination metadata
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	hasMore := page < totalPages

	pagination := response.Pagination{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
		HasMore:    hasMore,
	}

	response.Paginated(w, http.StatusOK, transfers, pagination, r.Context())
}

// GetByID retrieves a specific transfer by ID
// GET /api/transfers/{id}
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Extract transfer ID from URL
	transferID := chi.URLParam(r, "id")
	if transferID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Transfer ID is required", nil)
		return
	}

	// Get transfer
	transfer, err := h.service.GetByID(r.Context(), userID, transferID)
	if err != nil {
		if err == ErrTransferNotFound {
			response.Error(w, http.StatusNotFound, "RES_003", "Transfer not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "SYS_001", "Internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, transfer, r.Context())
}

// Cancel cancels a pending transfer
// POST /api/transfers/{id}/cancel
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Extract transfer ID from URL
	transferID := chi.URLParam(r, "id")
	if transferID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Transfer ID is required", nil)
		return
	}

	// Cancel transfer
	err := h.service.Cancel(r.Context(), userID, transferID)
	if err != nil {
		if err == ErrTransferNotFound {
			response.Error(w, http.StatusNotFound, "RES_003", "Transfer not found", nil)
			return
		}
		if err == ErrInvalidTransferStatus {
			response.Error(w, http.StatusBadRequest, "BUS_004", "Cannot cancel transfer with current status", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "SYS_001", "Internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "Transfer cancelled successfully"}, r.Context())
}

// handleTransferError maps service errors to HTTP responses
func (h *Handler) handleTransferError(w http.ResponseWriter, err error) {
	switch err {
	case ErrInsufficientBalance:
		response.Error(w, http.StatusBadRequest, "BUS_001", "Insufficient balance", nil)
	case ErrDailyLimitExceeded:
		response.Error(w, http.StatusBadRequest, "BUS_002", "Daily transfer limit exceeded", nil)
	case ErrMonthlyLimitExceeded:
		response.Error(w, http.StatusBadRequest, "BUS_003", "Monthly transfer limit exceeded", nil)
	case ErrInvalidAmount:
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid transfer amount", nil)
	case ErrInvalidPIXKey:
		response.Error(w, http.StatusBadRequest, "VAL_002", "Invalid PIX key", nil)
	case ErrInvalidBankData:
		response.Error(w, http.StatusBadRequest, "VAL_003", "Invalid bank account data", nil)
	case ErrInvalidCPF:
		response.Error(w, http.StatusBadRequest, "VAL_004", "Invalid CPF", nil)
	case ErrRecipientNotFound:
		response.Error(w, http.StatusNotFound, "RES_002", "Recipient user not found", nil)
	case ErrCannotTransferToSelf:
		response.Error(w, http.StatusBadRequest, "BUS_006", "Cannot transfer to yourself", nil)
	default:
		response.Error(w, http.StatusInternalServerError, "SYS_001", "Internal server error", nil)
	}
}
