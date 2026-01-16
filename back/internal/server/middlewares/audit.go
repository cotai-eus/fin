package middlewares

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/google/uuid"
	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
	"github.com/sqlc-dev/pqtype"
)

// AuditLogger handles audit logging for mutations
type AuditLogger struct {
	db      *sql.DB
	queries *db.Queries
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(database *sql.DB) *AuditLogger {
	return &AuditLogger{
		db:      database,
		queries: db.New(database),
	}
}

// AuditContext holds audit information to be logged
type AuditContext struct {
	ResourceType string
	ResourceID   uuid.UUID
	OldValues    map[string]interface{}
	NewValues    map[string]interface{}
}

// SetAuditContext adds audit context to the request context
func SetAuditContext(ctx context.Context, auditCtx *AuditContext) context.Context {
	return context.WithValue(ctx, "audit_context", auditCtx)
}

// GetAuditContext retrieves audit context from the request context
func GetAuditContext(ctx context.Context) *AuditContext {
	if auditCtx, ok := ctx.Value("audit_context").(*AuditContext); ok {
		return auditCtx
	}
	return nil
}

// responseWriter wraps http.ResponseWriter to capture status code
type auditResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func newAuditResponseWriter(w http.ResponseWriter) *auditResponseWriter {
	return &auditResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		body:           &bytes.Buffer{},
	}
}

func (w *auditResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *auditResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// AuditMiddleware creates middleware that logs mutations to audit_logs table
func (al *AuditLogger) AuditMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Only audit mutations (POST, PATCH, DELETE, PUT)
			if !isMutationMethod(r.Method) {
				next.ServeHTTP(w, r)
				return
			}

			// Capture request body for new_values
			var requestBody map[string]interface{}
			if r.Body != nil {
				bodyBytes, err := io.ReadAll(r.Body)
				if err == nil {
					// Restore body for downstream handlers
					r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

					// Try to parse as JSON
					if len(bodyBytes) > 0 {
						json.Unmarshal(bodyBytes, &requestBody)
					}
				}
			}

			// Wrap response writer to capture status
			arw := newAuditResponseWriter(w)

			// Execute the handler
			next.ServeHTTP(arw, r)

			// Log audit entry asynchronously (don't block response)
			go al.logAudit(r, arw, requestBody)
		})
	}
}

func (al *AuditLogger) logAudit(r *http.Request, arw *auditResponseWriter, requestBody map[string]interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Extract user ID from context
	var userID uuid.UUID
	if userIDStr, ok := r.Context().Value("user_id").(string); ok && userIDStr != "" {
		if parsed, err := uuid.Parse(userIDStr); err == nil {
			userID = parsed
		}
	}

	// Get audit context if set by handler
	auditCtx := GetAuditContext(r.Context())
	if auditCtx == nil {
		// Default audit context
		auditCtx = &AuditContext{
			ResourceType: inferResourceType(r.URL.Path),
			ResourceID:   uuid.Nil,
			NewValues:    requestBody,
		}
	}

	// Extract IP address
	ipAddress := extractIP(r)

	// Extract request ID
	requestID := r.Context().Value("request_id")
	var requestIDStr string
	if requestID != nil {
		requestIDStr = requestID.(string)
	}

	// Determine status (success/failure based on HTTP status)
	status := "success"
	if arw.statusCode >= 400 {
		status = "failure"
	}

	// Convert old_values and new_values to JSON
	var oldValuesJSON, newValuesJSON pqtype.NullRawMessage
	if auditCtx.OldValues != nil {
		if data, err := json.Marshal(auditCtx.OldValues); err == nil {
			oldValuesJSON = pqtype.NullRawMessage{RawMessage: data, Valid: true}
		}
	}
	if auditCtx.NewValues != nil {
		if data, err := json.Marshal(auditCtx.NewValues); err == nil {
			newValuesJSON = pqtype.NullRawMessage{RawMessage: data, Valid: true}
		}
	}

	// Parse IP address
	var ipAddr pqtype.Inet
	if parsedIP := net.ParseIP(ipAddress); parsedIP != nil {
		ipAddr = pqtype.Inet{IPNet: net.IPNet{IP: parsedIP}, Valid: true}
	}

	// Insert audit log
	_, err := al.queries.CreateAuditLog(ctx, db.CreateAuditLogParams{
		UserID:       uuid.NullUUID{UUID: userID, Valid: userID != uuid.Nil},
		Action:       r.Method + " " + r.URL.Path,
		ResourceType: auditCtx.ResourceType,
		ResourceID:   auditCtx.ResourceID,
		OldValues:    oldValuesJSON,
		NewValues:    newValuesJSON,
		IpAddress:    ipAddr,
		UserAgent:    sql.NullString{String: r.UserAgent(), Valid: r.UserAgent() != ""},
		RequestID:    sql.NullString{String: requestIDStr, Valid: requestIDStr != ""},
		Status:       status,
	})

	if err != nil {
		// Log error but don't fail the request
		// In production, send to monitoring/alerting system
		_ = err
	}
}

// isMutationMethod checks if HTTP method is a mutation
func isMutationMethod(method string) bool {
	return method == http.MethodPost ||
		method == http.MethodPut ||
		method == http.MethodPatch ||
		method == http.MethodDelete
}

// inferResourceType tries to infer resource type from URL path
func inferResourceType(path string) string {
	// Simple inference: /api/transfers -> TRANSFER
	// /api/cards -> CARD, etc.
	if len(path) > 5 && path[:5] == "/api/" {
		parts := splitPath(path[5:])
		if len(parts) > 0 {
			resource := parts[0]
			// Singularize and uppercase
			if len(resource) > 0 {
				if resource[len(resource)-1] == 's' {
					resource = resource[:len(resource)-1]
				}
				return toUpper(resource)
			}
		}
	}
	return "UNKNOWN"
}

// splitPath splits path by /
func splitPath(path string) []string {
	var parts []string
	var current string
	for _, c := range path {
		if c == '/' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

// toUpper converts string to uppercase
func toUpper(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			result[i] = c - 32
		} else {
			result[i] = c
		}
	}
	return string(result)
}

// extractIP extracts the real IP address from the request
func extractIP(r *http.Request) string {
	// Check X-Forwarded-For header first (from APISIX)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// Take the first IP in the list
		if ip, _, err := net.SplitHostPort(xff); err == nil {
			return ip
		}
		return xff
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}

	return r.RemoteAddr
}
