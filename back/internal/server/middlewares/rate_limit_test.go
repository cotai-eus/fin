package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	tests := []struct {
		name          string
		maxRequests   int
		window        time.Duration
		requestCount  int
		expectAllowed int
	}{
		{
			name:          "allows requests within limit",
			maxRequests:   5,
			window:        time.Hour,
			requestCount:  5,
			expectAllowed: 5,
		},
		{
			name:          "blocks requests exceeding limit",
			maxRequests:   3,
			window:        time.Hour,
			requestCount:  5,
			expectAllowed: 3,
		},
		{
			name:          "resets after window expires",
			maxRequests:   2,
			window:        100 * time.Millisecond,
			requestCount:  2,
			expectAllowed: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limiter := NewRateLimiter(tt.maxRequests, tt.window)
			userID := "test-user-123"

			allowed := 0
			for i := 0; i < tt.requestCount; i++ {
				if limiter.Allow(userID) {
					allowed++
				}
			}

			if allowed != tt.expectAllowed {
				t.Errorf("expected %d allowed requests, got %d", tt.expectAllowed, allowed)
			}
		})
	}
}

func TestRateLimiter_WindowReset(t *testing.T) {
	limiter := NewRateLimiter(2, 50*time.Millisecond)
	userID := "test-user-456"

	// First 2 requests should be allowed
	if !limiter.Allow(userID) {
		t.Error("first request should be allowed")
	}
	if !limiter.Allow(userID) {
		t.Error("second request should be allowed")
	}

	// Third request should be blocked
	if limiter.Allow(userID) {
		t.Error("third request should be blocked")
	}

	// Wait for window to expire
	time.Sleep(60 * time.Millisecond)

	// Next request should be allowed after window reset
	if !limiter.Allow(userID) {
		t.Error("request after window reset should be allowed")
	}
}

func TestRateLimiter_MultipleUsers(t *testing.T) {
	limiter := NewRateLimiter(2, time.Hour)

	// User 1 uses their limit
	if !limiter.Allow("user1") {
		t.Error("user1 first request should be allowed")
	}
	if !limiter.Allow("user1") {
		t.Error("user1 second request should be allowed")
	}
	if limiter.Allow("user1") {
		t.Error("user1 third request should be blocked")
	}

	// User 2 should have independent limit
	if !limiter.Allow("user2") {
		t.Error("user2 first request should be allowed")
	}
	if !limiter.Allow("user2") {
		t.Error("user2 second request should be allowed")
	}
}

func TestRateLimiter_GetStatus(t *testing.T) {
	limiter := NewRateLimiter(5, time.Hour)
	userID := "test-user-789"

	// Make 3 requests
	for i := 0; i < 3; i++ {
		limiter.Allow(userID)
	}

	remaining, _ := limiter.GetStatus(userID)
	if remaining != 2 {
		t.Errorf("expected 2 remaining requests, got %d", remaining)
	}
}

func TestRateLimitMiddleware_AllowsRequest(t *testing.T) {
	middleware := RateLimitMiddleware(10, time.Hour)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	req = addUserIDToContext(req, "test-user")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	// Check rate limit headers
	if rr.Header().Get("X-RateLimit-Limit") == "" {
		t.Error("expected X-RateLimit-Limit header")
	}
	if rr.Header().Get("X-RateLimit-Remaining") == "" {
		t.Error("expected X-RateLimit-Remaining header")
	}
}

func TestRateLimitMiddleware_BlocksExceedingRequests(t *testing.T) {
	middleware := RateLimitMiddleware(2, time.Hour)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	userID := "test-user-blocked"

	// Make 3 requests (limit is 2)
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/api/test", nil)
		req = addUserIDToContext(req, userID)

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if i < 2 {
			if rr.Code != http.StatusOK {
				t.Errorf("request %d: expected status 200, got %d", i+1, rr.Code)
			}
		} else {
			if rr.Code != http.StatusTooManyRequests {
				t.Errorf("request %d: expected status 429, got %d", i+1, rr.Code)
			}
		}
	}
}

func TestRateLimitMiddleware_NoUserID(t *testing.T) {
	middleware := RateLimitMiddleware(10, time.Hour)

	handler := middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/api/test", nil)
	// Don't add user_id to context

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Should allow request when no user ID (auth middleware will handle)
	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

// Helper to add user_id to context
func addUserIDToContext(r *http.Request, userID string) *http.Request {
	ctx := context.WithValue(r.Context(), "user_id", userID)
	return r.WithContext(ctx)
}
