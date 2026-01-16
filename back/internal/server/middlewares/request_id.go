package middlewares

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const (
	// HeaderRequestID is the header name for request ID
	HeaderRequestID = "X-Request-ID"
)

// RequestID middleware generates or extracts request ID for tracing
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Try to get existing request ID from header
		requestID := r.Header.Get(HeaderRequestID)

		// Generate new ID if not present
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Store in context
		ctx := context.WithValue(r.Context(), "request_id", requestID)

		// Add to response header
		w.Header().Set(HeaderRequestID, requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// generateRequestID creates a random 16-byte hex string
func generateRequestID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// Fallback to timestamp-based ID if crypto/rand fails
		return hex.EncodeToString([]byte("fallback"))
	}
	return hex.EncodeToString(b)
}
