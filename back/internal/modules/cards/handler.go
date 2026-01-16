package cards

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/lauratech/fin/back/internal/shared/response"
)

// Handler handles HTTP requests for cards
type Handler struct {
	service *Service
}

// NewHandler creates a new card handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateCard creates a new card
// POST /api/cards
func (h *Handler) CreateCard(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Decode request
	var req CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	// Create card
	card, err := h.service.CreateCard(r.Context(), userID, req)
	if err != nil {
		h.handleCardError(w, err)
		return
	}

	// Return card summary (not full card with sensitive data)
	summary := &CardSummary{
		ID:                       card.ID,
		Type:                     card.Type,
		Brand:                    card.Brand,
		Status:                   card.Status,
		LastFourDigits:           card.LastFourDigits,
		HolderName:               card.HolderName,
		ExpiryMonth:              card.ExpiryMonth,
		ExpiryYear:               card.ExpiryYear,
		DailyLimitCents:          card.DailyLimitCents,
		MonthlyLimitCents:        card.MonthlyLimitCents,
		CurrentDailySpentCents:   card.CurrentDailySpentCents,
		CurrentMonthlySpentCents: card.CurrentMonthlySpentCents,
		IsContactless:            card.IsContactless,
		IsInternational:          card.IsInternational,
		CreatedAt:                card.CreatedAt,
	}

	response.Success(w, http.StatusCreated, summary, r.Context())
}

// ListCards lists user's cards (without sensitive data)
// GET /api/cards
func (h *Handler) ListCards(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// List cards
	cards, err := h.service.ListUserCards(r.Context(), userID)
	if err != nil {
		h.handleCardError(w, err)
		return
	}

	response.Success(w, http.StatusOK, cards, r.Context())
}

// GetCardDetails retrieves full card details (with masked number)
// GET /api/cards/{id}
func (h *Handler) GetCardDetails(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Extract card ID from URL
	cardID := chi.URLParam(r, "id")
	if cardID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_002", "Card ID is required", nil)
		return
	}

	// Get card details
	details, err := h.service.GetCardByID(r.Context(), userID, cardID)
	if err != nil {
		h.handleCardError(w, err)
		return
	}

	response.Success(w, http.StatusOK, details, r.Context())
}

// BlockCard blocks a card
// POST /api/cards/{id}/block
func (h *Handler) BlockCard(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Extract card ID from URL
	cardID := chi.URLParam(r, "id")
	if cardID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_002", "Card ID is required", nil)
		return
	}

	// Block card
	if err := h.service.BlockCard(r.Context(), userID, cardID); err != nil {
		h.handleCardError(w, err)
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "Card blocked successfully"}, r.Context())
}

// UnblockCard unblocks a card
// POST /api/cards/{id}/unblock
func (h *Handler) UnblockCard(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Extract card ID from URL
	cardID := chi.URLParam(r, "id")
	if cardID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_002", "Card ID is required", nil)
		return
	}

	// Unblock card
	if err := h.service.UnblockCard(r.Context(), userID, cardID); err != nil {
		h.handleCardError(w, err)
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "Card unblocked successfully"}, r.Context())
}

// UpdateLimits updates spending limits
// PATCH /api/cards/{id}/limits
func (h *Handler) UpdateLimits(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Extract card ID from URL
	cardID := chi.URLParam(r, "id")
	if cardID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_002", "Card ID is required", nil)
		return
	}

	// Decode request
	var req UpdateLimitsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	// Update limits
	if err := h.service.UpdateLimits(r.Context(), userID, cardID, req); err != nil {
		h.handleCardError(w, err)
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "Limits updated successfully"}, r.Context())
}

// UpdateSecuritySettings updates security settings
// PATCH /api/cards/{id}/security
func (h *Handler) UpdateSecuritySettings(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Extract card ID from URL
	cardID := chi.URLParam(r, "id")
	if cardID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_002", "Card ID is required", nil)
		return
	}

	// Decode request
	var req SecuritySettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	// Update settings
	if err := h.service.UpdateSecuritySettings(r.Context(), userID, cardID, req); err != nil {
		h.handleCardError(w, err)
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "Security settings updated successfully"}, r.Context())
}

// SetPIN sets or changes card PIN
// POST /api/cards/{id}/pin
func (h *Handler) SetPIN(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Extract card ID from URL
	cardID := chi.URLParam(r, "id")
	if cardID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_002", "Card ID is required", nil)
		return
	}

	// Decode request
	var req SetPINRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	// Set PIN
	if err := h.service.SetPIN(r.Context(), userID, cardID, req); err != nil {
		h.handleCardError(w, err)
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "PIN set successfully"}, r.Context())
}

// CancelCard cancels a card
// DELETE /api/cards/{id}
func (h *Handler) CancelCard(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Extract card ID from URL
	cardID := chi.URLParam(r, "id")
	if cardID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_002", "Card ID is required", nil)
		return
	}

	// Decode request (optional reason)
	var req CancelCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// If no body provided, use default reason
		req.Reason = "user_request"
	}

	// Cancel card
	if err := h.service.CancelCard(r.Context(), userID, cardID, req.Reason); err != nil {
		h.handleCardError(w, err)
		return
	}

	response.Success(w, http.StatusOK, map[string]string{"message": "Card cancelled successfully"}, r.Context())
}

// handleCardError maps service errors to HTTP responses
func (h *Handler) handleCardError(w http.ResponseWriter, err error) {
	switch err {
	// Card errors
	case ErrCardNotFound:
		response.Error(w, http.StatusNotFound, "CARD_001", "Card not found", nil)
	case ErrCardNotActive:
		response.Error(w, http.StatusBadRequest, "CARD_002", "Card is not active", nil)
	case ErrCardExpired:
		response.Error(w, http.StatusBadRequest, "CARD_003", "Card has expired", nil)
	case ErrCardCancelled:
		response.Error(w, http.StatusBadRequest, "CARD_004", "Card has been cancelled", nil)
	case ErrCardBlocked:
		response.Error(w, http.StatusBadRequest, "CARD_005", "Card is blocked", nil)

	// Validation errors
	case ErrInvalidCardNumber:
		response.Error(w, http.StatusBadRequest, "VAL_003", "Invalid card number", nil)
	case ErrInvalidCVV:
		response.Error(w, http.StatusBadRequest, "VAL_004", "Invalid CVV", nil)
	case ErrInvalidPIN:
		response.Error(w, http.StatusBadRequest, "VAL_005", "Invalid PIN format (must be 4-6 digits)", nil)
	case ErrWeakPIN:
		response.Error(w, http.StatusBadRequest, "VAL_006", "Weak PIN detected (avoid sequences and repeating digits)", nil)
	case ErrInvalidExpiryDate:
		response.Error(w, http.StatusBadRequest, "VAL_007", "Invalid expiry date", nil)
	case ErrInvalidCardType:
		response.Error(w, http.StatusBadRequest, "VAL_008", "Invalid card type (must be physical or virtual)", nil)
	case ErrInvalidCardBrand:
		response.Error(w, http.StatusBadRequest, "VAL_009", "Invalid card brand (must be visa, mastercard, or elo)", nil)

	// Limit errors
	case ErrDailyLimitExceeded:
		response.Error(w, http.StatusBadRequest, "LIMIT_001", "Daily spending limit exceeded", nil)
	case ErrMonthlyLimitExceeded:
		response.Error(w, http.StatusBadRequest, "LIMIT_002", "Monthly spending limit exceeded", nil)
	case ErrInvalidLimit:
		response.Error(w, http.StatusBadRequest, "LIMIT_003", "Invalid limit amount", nil)

	// PIN errors
	case ErrPINNotSet:
		response.Error(w, http.StatusBadRequest, "PIN_001", "PIN not set for this card", nil)
	case ErrPINIncorrect:
		response.Error(w, http.StatusUnauthorized, "PIN_002", "Incorrect PIN", nil)
	case ErrPINMismatch:
		response.Error(w, http.StatusUnauthorized, "PIN_003", "Current PIN does not match", nil)

	// Security errors
	case ErrInternationalBlocked:
		response.Error(w, http.StatusForbidden, "SEC_001", "International transactions are blocked", nil)
	case ErrOnlineBlocked:
		response.Error(w, http.StatusForbidden, "SEC_002", "Online transactions are blocked", nil)
	case ErrContactlessBlocked:
		response.Error(w, http.StatusForbidden, "SEC_003", "Contactless transactions are blocked", nil)

	// Encryption errors
	case ErrEncryptionFailed:
		response.Error(w, http.StatusInternalServerError, "SYS_002", "Encryption failed", nil)
	case ErrDecryptionFailed:
		response.Error(w, http.StatusInternalServerError, "SYS_003", "Decryption failed", nil)

	// Authorization errors
	case ErrUnauthorized:
		response.Error(w, http.StatusForbidden, "AUTH_002", "Unauthorized access to card", nil)

	// Default error
	default:
		response.Error(w, http.StatusInternalServerError, "SYS_001", "Internal server error", nil)
	}
}
