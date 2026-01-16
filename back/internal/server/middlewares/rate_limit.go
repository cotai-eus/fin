package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/lauratech/fin/back/internal/shared/response"
)

// RateLimiter tracks request counts per user within time windows
type RateLimiter struct {
	mu       sync.RWMutex
	limits   map[string]*userLimit // key: userID
	max      int                   // maximum requests allowed
	window   time.Duration         // time window for rate limit
	cleanTTL time.Duration         // how long to keep inactive user entries
}

type userLimit struct {
	count       int
	windowStart time.Time
	lastAccess  time.Time
}

// NewRateLimiter creates a new rate limiter with specified limits
func NewRateLimiter(maxRequests int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		limits:   make(map[string]*userLimit),
		max:      maxRequests,
		window:   window,
		cleanTTL: window * 2, // Keep entries for 2x window duration
	}

	// Start background cleanup goroutine
	go rl.cleanupLoop()

	return rl
}

// Allow checks if a request from the given user should be allowed
func (rl *RateLimiter) Allow(userID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()

	limit, exists := rl.limits[userID]
	if !exists {
		rl.limits[userID] = &userLimit{
			count:       1,
			windowStart: now,
			lastAccess:  now,
		}
		return true
	}

	// Check if we need to reset the window
	if now.Sub(limit.windowStart) >= rl.window {
		limit.count = 1
		limit.windowStart = now
		limit.lastAccess = now
		return true
	}

	// Within the window, check if limit exceeded
	limit.lastAccess = now
	if limit.count >= rl.max {
		return false
	}

	limit.count++
	return true
}

// GetStatus returns current rate limit status for a user
func (rl *RateLimiter) GetStatus(userID string) (remaining int, resetAt time.Time) {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	limit, exists := rl.limits[userID]
	if !exists {
		return rl.max, time.Now()
	}

	resetAt = limit.windowStart.Add(rl.window)
	remaining = rl.max - limit.count
	if remaining < 0 {
		remaining = 0
	}

	return remaining, resetAt
}

// cleanupLoop periodically removes stale entries to prevent memory leaks
func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(rl.window)
	defer ticker.Stop()

	for range ticker.C {
		rl.cleanup()
	}
}

func (rl *RateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	for userID, limit := range rl.limits {
		if now.Sub(limit.lastAccess) > rl.cleanTTL {
			delete(rl.limits, userID)
		}
	}
}

// RateLimitMiddleware creates a middleware that enforces rate limits per user
func RateLimitMiddleware(maxRequests int, window time.Duration) func(http.Handler) http.Handler {
	limiter := NewRateLimiter(maxRequests, window)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract user ID from context (set by auth middleware)
			userID, ok := r.Context().Value("user_id").(string)
			if !ok || userID == "" {
				// If no user ID, allow the request (auth middleware will handle)
				next.ServeHTTP(w, r)
				return
			}

			// Check rate limit
			if !limiter.Allow(userID) {
				remaining, resetAt := limiter.GetStatus(userID)

				// Add rate limit headers
				w.Header().Set("X-RateLimit-Limit", formatInt(maxRequests))
				w.Header().Set("X-RateLimit-Remaining", formatInt(remaining))
				w.Header().Set("X-RateLimit-Reset", formatInt(int(resetAt.Unix())))
				w.Header().Set("Retry-After", formatInt(int(time.Until(resetAt).Seconds())))

				response.Error(w, http.StatusTooManyRequests, "SYS_003", "Rate limit exceeded", map[string]interface{}{
					"limit":       maxRequests,
					"window":      window.String(),
					"reset_at":    resetAt.Format(time.RFC3339),
					"retry_after": int(time.Until(resetAt).Seconds()),
				})
				return
			}

			// Add rate limit info headers for successful requests
			remaining, resetAt := limiter.GetStatus(userID)
			w.Header().Set("X-RateLimit-Limit", formatInt(maxRequests))
			w.Header().Set("X-RateLimit-Remaining", formatInt(remaining))
			w.Header().Set("X-RateLimit-Reset", formatInt(int(resetAt.Unix())))

			next.ServeHTTP(w, r)
		})
	}
}

// Helper to format int as string
func formatInt(n int) string {
	return fmt.Sprintf("%d", n)
}
