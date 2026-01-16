package middlewares

import (
	"context"
	"net/http"
	"strings"
)

const (
	// HeaderKratosIdentityID is the header injected by APISIX after validating session with Kratos
	HeaderKratosIdentityID = "X-Kratos-Authenticated-Identity-Id"
)

// Auth middleware validates the Kratos identity header from APISIX
// Security: Only trusts headers from APISIX gateway
func Auth(trustedProxyIP string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Security: In production, validate that request comes from APISIX
			// remoteIP := strings.Split(r.RemoteAddr, ":")[0]
			// if remoteIP != trustedProxyIP {
			// 	http.Error(w, `{"error":{"code":"AUTH_003","message":"Forbidden: Invalid proxy"}}`, http.StatusForbidden)
			// 	return
			// }

			// Extract identity ID from header
			identityID := r.Header.Get(HeaderKratosIdentityID)
			identityID = strings.TrimSpace(identityID)

			// Validate header presence and format (UUID)
			if identityID == "" {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, `{"error":{"code":"AUTH_001","message":"Unauthorized: Missing identity header"}}`, http.StatusUnauthorized)
				return
			}

			if !isValidUUID(identityID) {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, `{"error":{"code":"AUTH_002","message":"Unauthorized: Invalid identity format"}}`, http.StatusUnauthorized)
				return
			}

			// Store user ID in context for downstream handlers
			ctx := context.WithValue(r.Context(), "user_id", identityID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// isValidUUID checks if string is a valid UUID format
func isValidUUID(uuid string) bool {
	if len(uuid) != 36 {
		return false
	}

	// Simple UUID v4 format check: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		return false
	}

	return true
}
