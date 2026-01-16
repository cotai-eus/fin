package response

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Data interface{} `json:"data"`
	Meta Meta        `json:"meta"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
	Meta  Meta        `json:"meta"`
}

// ErrorDetail contains error information
type ErrorDetail struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// Meta contains metadata about the response
type Meta struct {
	RequestID string `json:"request_id,omitempty"`
	Timestamp string `json:"timestamp"`
}

// Success writes a successful JSON response
func Success(w http.ResponseWriter, status int, data interface{}, ctx context.Context) {
	response := SuccessResponse{
		Data: data,
		Meta: Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Extract request ID from context
	if requestID, ok := ctx.Value("request_id").(string); ok {
		response.Meta.RequestID = requestID
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// Error writes an error JSON response
func Error(w http.ResponseWriter, status int, code, message string, details interface{}) {
	response := ErrorResponse{
		Error: ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
		Meta: Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
	Meta       Meta        `json:"meta"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	Total      int  `json:"total"`
	TotalPages int  `json:"total_pages"`
	HasMore    bool `json:"has_more"`
}

// Paginated writes a paginated JSON response
func Paginated(w http.ResponseWriter, status int, data interface{}, pagination Pagination, ctx context.Context) {
	response := PaginatedResponse{
		Data:       data,
		Pagination: pagination,
		Meta: Meta{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}

	// Extract request ID from context
	if requestID, ok := ctx.Value("request_id").(string); ok {
		response.Meta.RequestID = requestID
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}
