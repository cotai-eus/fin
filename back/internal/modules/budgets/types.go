package budgets

import "time"

// Budget represents a budget in the system
type Budget struct {
	ID                string    `json:"id"`
	UserID            string    `json:"user_id"`
	Category          string    `json:"category"`
	Period            string    `json:"period"` // "weekly", "monthly", "annual"
	LimitCents        int64     `json:"limit_cents"`
	CurrentSpentCents int64     `json:"current_spent_cents"`
	AlertThreshold    int16     `json:"alert_threshold"` // Percentage (0-100)
	AlertsEnabled     bool      `json:"alerts_enabled"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	CreatedAt         time.Time `json:"created_at"`
}

// BudgetSummary represents a budget summary for dashboard
type BudgetSummary struct {
	ID                string  `json:"id"`
	Category          string  `json:"category"`
	Period            string  `json:"period"`
	LimitCents        int64   `json:"limit_cents"`
	CurrentSpentCents int64   `json:"current_spent_cents"`
	PercentageUsed    float64 `json:"percentage_used"`
	AlertThreshold    int16   `json:"alert_threshold"`
	IsOverBudget      bool    `json:"is_over_budget"`
	IsNearLimit       bool    `json:"is_near_limit"`
}

// CreateBudgetRequest represents a request to create a budget
type CreateBudgetRequest struct {
	Category       string `json:"category"`
	Period         string `json:"period"` // "weekly", "monthly", "annual"
	LimitCents     int64  `json:"limit_cents"`
	AlertThreshold int16  `json:"alert_threshold,omitempty"` // Default: 75
	AlertsEnabled  bool   `json:"alerts_enabled"`            // Default: true
}

// UpdateBudgetRequest represents a request to update a budget
type UpdateBudgetRequest struct {
	LimitCents     int64 `json:"limit_cents,omitempty"`
	AlertThreshold int16 `json:"alert_threshold,omitempty"`
	AlertsEnabled  *bool `json:"alerts_enabled,omitempty"`
}

// BudgetListParams represents pagination parameters for listing budgets
type BudgetListParams struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Category string `json:"category,omitempty"` // Filter by category
	Period   string `json:"period,omitempty"`   // Filter by period
}

// CategorySpending represents spending by category
type CategorySpending struct {
	Category    string  `json:"category"`
	TotalCents  int64   `json:"total_cents"`
	Percentage  float64 `json:"percentage"`
	BudgetCents int64   `json:"budget_cents,omitempty"`
}

// SpendingTrend represents spending trend over time
type SpendingTrend struct {
	Date        string `json:"date"`
	AmountCents int64  `json:"amount_cents"`
}

// BudgetAnalytics represents budget analytics data
type BudgetAnalytics struct {
	TotalBudgetCents     int64               `json:"total_budget_cents"`
	TotalSpentCents      int64               `json:"total_spent_cents"`
	PercentageUsed       float64             `json:"percentage_used"`
	CategoriesOverBudget int                 `json:"categories_over_budget"`
	CategoryBreakdown    []*CategorySpending `json:"category_breakdown"`
}
