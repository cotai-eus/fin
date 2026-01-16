package budgets

import "errors"

var (
	// ErrBudgetNotFound is returned when a budget is not found
	ErrBudgetNotFound = errors.New("budget not found")

	// ErrInvalidCategory is returned when budget category is invalid
	ErrInvalidCategory = errors.New("invalid budget category")

	// ErrInvalidPeriod is returned when budget period is invalid
	ErrInvalidPeriod = errors.New("invalid budget period")

	// ErrInvalidLimit is returned when budget limit is invalid
	ErrInvalidLimit = errors.New("invalid budget limit")

	// ErrInvalidThreshold is returned when alert threshold is invalid
	ErrInvalidThreshold = errors.New("invalid alert threshold")

	// ErrInvalidDateRange is returned when date range is invalid
	ErrInvalidDateRange = errors.New("invalid date range")

	// ErrBudgetAlreadyExists is returned when a budget for category/period already exists
	ErrBudgetAlreadyExists = errors.New("budget already exists for this category and period")

	// ErrUnauthorized is returned when user doesn't own the budget
	ErrUnauthorized = errors.New("unauthorized to access this budget")
)
