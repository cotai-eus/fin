package middlewares

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += n
	return n, err
}

// Logger middleware logs HTTP requests with structured format
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Process request
		next.ServeHTTP(wrapped, r)

		// Log request
		duration := time.Since(start)
		requestID := r.Context().Value("request_id")
		userID := r.Context().Value("user_id")

		log.Printf(
			"[HTTP] %s %s | status=%d | duration=%v | bytes=%d | ip=%s | request_id=%v | user_id=%v",
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
			wrapped.written,
			r.RemoteAddr,
			requestID,
			userID,
		)
	})
}
