package bills

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/lauratech/fin/back/internal/shared/response"
)

// Handler handles HTTP requests for bills
type Handler struct {
	service *Service
}

// NewHandler creates a new bill handler
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ValidateBarcode validates a bill barcode
// POST /api/bills/validate
func (h *Handler) ValidateBarcode(w http.ResponseWriter, r *http.Request) {
	var req ValidateBarcodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	result, err := h.service.ValidateBarcodeInfo(r.Context(), req.Barcode)
	if err != nil {
		h.handleBillError(w, err)
		return
	}

	response.Success(w, http.StatusOK, result, r.Context())
}

// CreateBill creates a new bill from a barcode
// POST /api/bills
func (h *Handler) CreateBill(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	var req CreateBillRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Invalid request body", nil)
		return
	}

	bill, err := h.service.CreateBill(r.Context(), userID, req)
	if err != nil {
		h.handleBillError(w, err)
		return
	}

	response.Success(w, http.StatusCreated, bill, r.Context())
}

// GetBill retrieves a bill by ID
// GET /api/bills/{id}
func (h *Handler) GetBill(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	billID := chi.URLParam(r, "id")
	if billID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Bill ID is required", nil)
		return
	}

	bill, err := h.service.GetBillByID(r.Context(), userID, billID)
	if err != nil {
		h.handleBillError(w, err)
		return
	}

	response.Success(w, http.StatusOK, bill, r.Context())
}

// ListBills lists bills for a user
// GET /api/bills
func (h *Handler) ListBills(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	// Parse query parameters
	params := BillListParams{
		Page:   1,
		Limit:  20,
		Status: r.URL.Query().Get("status"),
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

	bills, total, err := h.service.ListUserBills(r.Context(), userID, params)
	if err != nil {
		h.handleBillError(w, err)
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

	response.Paginated(w, http.StatusOK, bills, pagination, r.Context())
}

// PayBill processes a bill payment
// POST /api/bills/{id}/pay
func (h *Handler) PayBill(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	billID := chi.URLParam(r, "id")
	if billID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Bill ID is required", nil)
		return
	}

	bill, err := h.service.PayBill(r.Context(), userID, billID)
	if err != nil {
		h.handleBillError(w, err)
		return
	}

	response.Success(w, http.StatusOK, bill, r.Context())
}

// CancelBill cancels a bill
// DELETE /api/bills/{id}
func (h *Handler) CancelBill(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		response.Error(w, http.StatusUnauthorized, "AUTH_001", "Unauthorized", nil)
		return
	}

	billID := chi.URLParam(r, "id")
	if billID == "" {
		response.Error(w, http.StatusBadRequest, "VAL_001", "Bill ID is required", nil)
		return
	}

	err := h.service.CancelBill(r.Context(), userID, billID)
	if err != nil {
		h.handleBillError(w, err)
		return
	}

	response.Success(w, http.StatusNoContent, nil, r.Context())
}

// handleBillError maps domain errors to HTTP responses
func (h *Handler) handleBillError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrBillNotFound):
		response.Error(w, http.StatusNotFound, "BILL_001", "Bill not found", nil)
	case errors.Is(err, ErrInvalidBarcode):
		response.Error(w, http.StatusBadRequest, "BILL_002", "Invalid barcode format", nil)
	case errors.Is(err, ErrBarcodeAlreadyExists):
		response.Error(w, http.StatusConflict, "BILL_003", "Barcode already registered", nil)
	case errors.Is(err, ErrBillAlreadyPaid):
		response.Error(w, http.StatusBadRequest, "BILL_004", "Bill already paid", nil)
	case errors.Is(err, ErrBillCancelled):
		response.Error(w, http.StatusBadRequest, "BILL_005", "Bill is cancelled", nil)
	case errors.Is(err, ErrInsufficientBalance):
		response.Error(w, http.StatusPaymentRequired, "BILL_006", "Insufficient balance", nil)
	case errors.Is(err, ErrInvalidAmount):
		response.Error(w, http.StatusBadRequest, "BILL_007", "Invalid bill amount", nil)
	case errors.Is(err, ErrInvalidType):
		response.Error(w, http.StatusBadRequest, "BILL_008", "Invalid bill type", nil)
	case errors.Is(err, ErrUnauthorized):
		response.Error(w, http.StatusForbidden, "BILL_009", "Unauthorized to access this bill", nil)
	default:
		response.Error(w, http.StatusInternalServerError, "SYS_001", "Internal server error", nil)
	}
}
