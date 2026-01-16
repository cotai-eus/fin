package users

import (
	"encoding/json"
	"net/http"

	"github.com/lauratech/fin/back/internal/shared/response"
)

// Handler handles HTTP requests for users
type Handler struct {
	service *Service
}

// NewHandler creates a new user handler
func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

// GetCurrentUser retrieves the current authenticated user
// GET /api/users/me
func (h *Handler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	user, err := h.service.GetCurrentUser(r.Context(), userID)
	if err != nil {
		if err == ErrUserNotFound {
			response.Error(w, http.StatusNotFound, "USER_001", "User not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "SYS_001", "Internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, user, r.Context())
}

// CreateUser creates a new user
// POST /api/users
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	user, err := h.service.Create(r.Context(), req)
	if err != nil {
		if err == ErrUserAlreadyExists {
			response.Error(w, http.StatusConflict, "USER_002", "User already exists", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "SYS_001", "Internal server error", nil)
		return
	}

	response.Success(w, http.StatusCreated, user, r.Context())
}

// UpdateUser updates a user
// PATCH /api/users/{id}
func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL (to be implemented with chi router params)
	// For now, we'll use the authenticated user ID
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	user, err := h.service.Update(r.Context(), userID, req)
	if err != nil {
		if err == ErrUserNotFound {
			response.Error(w, http.StatusNotFound, "USER_001", "User not found", nil)
			return
		}
		response.Error(w, http.StatusInternalServerError, "SYS_001", "Internal server error", nil)
		return
	}

	response.Success(w, http.StatusOK, user, r.Context())
}
