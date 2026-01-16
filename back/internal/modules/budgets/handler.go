package budgets

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lauratech/fin/back/internal/shared/response"
)

// Handler handles HTTP requests for budgets
type Handler struct {
	service *Service
}

// NewHandler creates a new budget handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// CreateBudget creates a new budget
// POST /api/budgets
func (h *Handler) CreateBudget(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	var req CreateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	budget, err := h.service.CreateBudget(r.Context(), userID, req)
	if err != nil {
		h.handleBudgetError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, budget, r.Context())
}

// GetBudget retrieves a budget by ID
// GET /api/budgets/{id}
func (h *Handler) GetBudget(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	budgetID := chi.URLParam(r, "id")
	if budgetID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Budget ID is required", nil)
		return
	}

	budget, err := h.service.GetBudgetByID(r.Context(), userID, budgetID)
	if err != nil {
		h.handleBudgetError(w, err)
		return
	}

	response.Success(w, http.StatusOK, budget, r.Context())
}

// ListBudgets lists budgets for a user
// GET /api/budgets
func (h *Handler) ListBudgets(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Parse query parameters
	params := BudgetListParams{
		Page:     1,
		Limit:    20,
		Category: r.URL.Query().Get("category"),
		Period:   r.URL.Query().Get("period"),
	}

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			params.Page = page
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			params.Limit = limit
		}
	}

	budgets, total, err := h.service.ListUserBudgets(r.Context(), userID, params)
	if err != nil {
		h.handleBudgetError(w, err)
		return
	}

	// Calculate pagination
	totalPages := int(total) / params.Limit
	if int(total)%params.Limit != 0 {
		totalPages++
	}

	pagination := response.Pagination{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      int(total),
		TotalPages: totalPages,
		HasMore:    params.Page < totalPages,
	}

	response.Paginated(w, http.StatusOK, budgets, pagination, r.Context())
}

// UpdateBudget updates a budget
// PATCH /api/budgets/{id}
func (h *Handler) UpdateBudget(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	budgetID := chi.URLParam(r, "id")
	if budgetID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Budget ID is required", nil)
		return
	}

	var req UpdateBudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	budget, err := h.service.UpdateBudget(r.Context(), userID, budgetID, req)
	if err != nil {
		h.handleBudgetError(w, err)
		return
	}

	response.Success(w, http.StatusOK, budget, r.Context())
}

// DeleteBudget deletes a budget
// DELETE /api/budgets/{id}
func (h *Handler) DeleteBudget(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	budgetID := chi.URLParam(r, "id")
	if budgetID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Budget ID is required", nil)
		return
	}

	err := h.service.DeleteBudget(r.Context(), userID, budgetID)
	if err != nil {
		h.handleBudgetError(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, nil, r.Context())
}

// GetBudgetSummary retrieves budget summary/analytics
// GET /api/budgets/summary
func (h *Handler) GetBudgetSummary(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	summary, err := h.service.GetBudgetSummary(r.Context(), userID)
	if err != nil {
		h.handleBudgetError(w, err)
		return
	}

	response.Success(w, http.StatusOK, summary, r.Context())
}

// handleBudgetError maps domain errors to HTTP responses
func (h *Handler) handleBudgetError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrBudgetNotFound):
		response.Error(w, http.StatusNotFound, "BUDGET_001", "Budget not found", nil)
	case errors.Is(err, ErrInvalidCategory):
		response.Error(w, http.StatusBadRequest, "BUDGET_002", "Invalid budget category", nil)
	case errors.Is(err, ErrInvalidPeriod):
		response.Error(w, http.StatusBadRequest, "BUDGET_003", "Invalid budget period", nil)
	case errors.Is(err, ErrInvalidLimit):
		response.Error(w, http.StatusBadRequest, "BUDGET_004", "Invalid budget limit", nil)
	case errors.Is(err, ErrInvalidThreshold):
		response.Error(w, http.StatusBadRequest, "BUDGET_005", "Invalid alert threshold", nil)
	case errors.Is(err, ErrBudgetAlreadyExists):
		response.Error(w, http.StatusConflict, "BUDGET_006", "Budget already exists for this category and period", nil)
	case errors.Is(err, ErrUnauthorized):
		response.Error(w, http.StatusForbidden, "BUDGET_007", "Unauthorized to access this budget", nil)
	default:
		response.Error(w, http.StatusInternalServerError, "SYS_001", "Internal server error", nil)
	}
}

// GetCategorySpending returns spending breakdown by category
// GET /api/analytics/category-spending
func (h *Handler) GetCategorySpending(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Parse date range from query params (optional)
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")

	analytics, err := h.service.GetCategorySpending(r.Context(), userID, startDate, endDate)
	if err != nil {
		h.handleBudgetError(w, err)
		return
	}

	response.Success(w, http.StatusOK, analytics, r.Context())
}

// GetSpendingTrends returns spending trends over time
// GET /api/analytics/spending-trends
func (h *Handler) GetSpendingTrends(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	trends, err := h.service.GetSpendingTrends(r.Context(), userID)
	if err != nil {
		h.handleBudgetError(w, err)
		return
	}

	response.Success(w, http.StatusOK, trends, r.Context())
}
