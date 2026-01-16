package budgets

import (
	"time"
)

// ValidCategory represents valid budget categories
var ValidCategories = map[string]bool{
	"groceries":      true,
	"dining":         true,
	"transportation": true,
	"entertainment":  true,
	"shopping":       true,
	"bills":          true,
	"health":         true,
	"education":      true,
	"travel":         true,
	"other":          true,
}

// ValidPeriods represents valid budget periods
var ValidPeriods = map[string]bool{
	"weekly":  true,
	"monthly": true,
	"annual":  true,
}

// ValidateCategory validates budget category
func ValidateCategory(category string) error {
	if !ValidCategories[category] {
		return ErrInvalidCategory
	}
	return nil
}

// ValidatePeriod validates budget period
func ValidatePeriod(period string) error {
	if !ValidPeriods[period] {
		return ErrInvalidPeriod
	}
	return nil
}

// ValidateLimit validates budget limit
func ValidateLimit(limitCents int64) error {
	if limitCents <= 0 {
		return ErrInvalidLimit
	}

	// Maximum budget: R$ 1,000,000.00 (100,000,000 cents)
	if limitCents > 100000000 {
		return ErrInvalidLimit
	}

	return nil
}

// ValidateThreshold validates alert threshold
func ValidateThreshold(threshold int16) error {
	if threshold < 0 || threshold > 100 {
		return ErrInvalidThreshold
	}
	return nil
}

// CalculateDateRange calculates start and end dates for a period
func CalculateDateRange(period string) (time.Time, time.Time, error) {
	now := time.Now()

	switch period {
	case "weekly":
		// Start: Beginning of current week (Monday)
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // Sunday
		}
		start := now.AddDate(0, 0, -(weekday - 1))
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
		end := start.AddDate(0, 0, 7)
		return start, end, nil

	case "monthly":
		// Start: First day of current month
		start := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, 0)
		return start, end, nil

	case "annual":
		// Start: First day of current year
		start := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(1, 0, 0)
		return start, end, nil

	default:
		return time.Time{}, time.Time{}, ErrInvalidPeriod
	}
}
