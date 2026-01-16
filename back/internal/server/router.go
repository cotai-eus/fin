package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lauratech/fin/back/internal/server/middlewares"
)

// setupRouter configures all routes and middleware
func (s *Server) setupRouter() *chi.Mux {
	r := chi.NewRouter()

	// ========================================
	// Global Middleware Stack (Order Matters)
	// ========================================

	// 1. Recovery (must be first)
	r.Use(middleware.Recoverer)

	// 2. Request ID (for tracing)
	r.Use(middlewares.RequestID)

	// 3. Logger (after request ID)
	r.Use(middlewares.Logger)

	// 4. Timeout (30 seconds)
	r.Use(middleware.Timeout(30 * time.Second))

	// 5. CORS
	r.Use(middlewares.CORS)

	// ========================================
	// Public Routes (No Auth)
	// ========================================

	r.Get("/health", s.handleHealth)

	// ========================================
	// Protected Routes (Require Auth)
	// ========================================

	r.Group(func(r chi.Router) {
		// Authentication middleware (validates APISIX header)
		r.Use(middlewares.Auth(s.Config.TrustedProxyIP))

		// Audit middleware (logs mutations to audit_logs)
		auditLogger := middlewares.NewAuditLogger(s.DB)
		r.Use(auditLogger.AuditMiddleware())

		// API routes
		r.Route("/api", func(r chi.Router) {
			// Users
			r.Route("/users", func(r chi.Router) {
				r.Get("/me", s.usersHandler.GetCurrentUser)
				r.Post("/", s.usersHandler.CreateUser)
				r.Patch("/me", s.usersHandler.UpdateUser)
			})

			// Transfers (10 requests/hour per endpoint)
			r.Route("/transfers", func(r chi.Router) {
				r.With(middlewares.RateLimitMiddleware(10, time.Hour)).Post("/pix", s.transfersHandler.ExecutePIX)
				r.With(middlewares.RateLimitMiddleware(10, time.Hour)).Post("/ted", s.transfersHandler.ExecuteTED)
				r.With(middlewares.RateLimitMiddleware(10, time.Hour)).Post("/p2p", s.transfersHandler.ExecuteP2P)
				r.Get("/", s.transfersHandler.List)
				r.Get("/{id}", s.transfersHandler.GetByID)
				r.With(middlewares.RateLimitMiddleware(10, time.Hour)).Post("/{id}/cancel", s.transfersHandler.Cancel)
			})

			// Deposits
			r.With(middlewares.RateLimitMiddleware(10, time.Hour)).Post("/deposits", s.transfersHandler.ExecuteDeposit)

			// Payment Requests
			r.Route("/payment-requests", func(r chi.Router) {
				r.With(middlewares.RateLimitMiddleware(10, time.Hour)).Post("/", s.transfersHandler.CreatePaymentRequest)
				r.Get("/", s.transfersHandler.ListPaymentRequests)
				r.Post("/{id}/approve", s.transfersHandler.ApprovePaymentRequest)
				r.Post("/{id}/reject", s.transfersHandler.RejectPaymentRequest)
			})

			// Cards
			r.Route("/cards", func(r chi.Router) {
				r.With(middlewares.RateLimitMiddleware(5, time.Hour)).Post("/", s.cardsHandler.CreateCard)                            // 5/hour - virtual card creation
				r.With(middlewares.RateLimitMiddleware(100, time.Hour)).Get("/", s.cardsHandler.ListCards)                            // 100/hour
				r.With(middlewares.RateLimitMiddleware(10, time.Hour)).Get("/{id}", s.cardsHandler.GetCardDetails)                    // 10/hour - sensitive data
				r.Get("/{id}/transactions", s.cardsHandler.ListCardTransactions)                                                      // Card transactions
				r.With(middlewares.RateLimitMiddleware(20, time.Hour)).Post("/{id}/block", s.cardsHandler.BlockCard)                  // 20/hour
				r.With(middlewares.RateLimitMiddleware(20, time.Hour)).Post("/{id}/unblock", s.cardsHandler.UnblockCard)              // 20/hour
				r.With(middlewares.RateLimitMiddleware(20, time.Hour)).Patch("/{id}/limits", s.cardsHandler.UpdateLimits)             // 20/hour
				r.With(middlewares.RateLimitMiddleware(20, time.Hour)).Patch("/{id}/security", s.cardsHandler.UpdateSecuritySettings) // 20/hour
				r.With(middlewares.RateLimitMiddleware(3, time.Hour)).Post("/{id}/pin", s.cardsHandler.SetPIN)                        // 3/hour - very sensitive
				r.With(middlewares.RateLimitMiddleware(5, time.Hour)).Delete("/{id}", s.cardsHandler.CancelCard)                      // 5/hour
			})

			// Transactions (card transactions)
			r.Route("/transactions", func(r chi.Router) {
				r.Get("/", s.cardsHandler.ListUserTransactions)
				r.Post("/export", s.cardsHandler.ExportTransactions)
			})

			// Bills
			r.Route("/bills", func(r chi.Router) {
				r.With(middlewares.RateLimitMiddleware(20, time.Hour)).Post("/validate", s.billsHandler.ValidateBarcode) // 20/hour
				r.Post("/", s.billsHandler.CreateBill)
				r.Get("/", s.billsHandler.ListBills)
				r.Get("/{id}", s.billsHandler.GetBill)
				r.With(middlewares.RateLimitMiddleware(10, time.Hour)).Post("/{id}/pay", s.billsHandler.PayBill) // 10/hour
				r.Delete("/{id}", s.billsHandler.CancelBill)
			})

			// Budgets
			r.Route("/budgets", func(r chi.Router) {
				r.With(middlewares.RateLimitMiddleware(20, time.Hour)).Post("/", s.budgetsHandler.CreateBudget) // 20/hour
				r.Get("/", s.budgetsHandler.ListBudgets)
				r.Get("/summary", s.budgetsHandler.GetBudgetSummary)
				r.Get("/{id}", s.budgetsHandler.GetBudget)
				r.Patch("/{id}", s.budgetsHandler.UpdateBudget)
				r.Delete("/{id}", s.budgetsHandler.DeleteBudget)
			})

			// Analytics
			r.Route("/analytics", func(r chi.Router) {
				r.Get("/category-spending", s.budgetsHandler.GetCategorySpending)
				r.Get("/spending-trends", s.budgetsHandler.GetSpendingTrends)
			})

			// Support
			r.Route("/support", func(r chi.Router) {
				r.With(middlewares.RateLimitMiddleware(10, time.Hour)).Post("/tickets", s.supportHandler.CreateTicket) // 10/hour
				r.Get("/tickets", s.supportHandler.ListTickets)
				r.Get("/tickets/{id}", s.supportHandler.GetTicket)
				r.Get("/tickets/{id}/messages", s.supportHandler.GetTicketWithMessages)
				r.Post("/tickets/{id}/messages", s.supportHandler.AddMessage)
				r.Patch("/tickets/{id}/status", s.supportHandler.UpdateTicketStatus)
			})
		})
	})

	return r
}

// ========================================
// Handler: Health Check
// ========================================

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	// Check database connectivity
	if err := s.DB.PingContext(r.Context()); err != nil {
		http.Error(w, `{"status":"unhealthy","database":"disconnected"}`, http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","database":"connected"}`))
}
