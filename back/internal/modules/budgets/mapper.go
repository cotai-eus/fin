package budgets

import (
	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

// dbBudgetToBudget converts a database budget to domain budget
func dbBudgetToBudget(dbBudget *db.Budget) *Budget {
	currentSpentCents := int64(0)
	if dbBudget.CurrentSpentCents.Valid {
		currentSpentCents = dbBudget.CurrentSpentCents.Int64
	}

	alertThreshold := int16(75) // Default
	if dbBudget.AlertThreshold.Valid {
		alertThreshold = dbBudget.AlertThreshold.Int16
	}

	alertsEnabled := true // Default
	if dbBudget.AlertsEnabled.Valid {
		alertsEnabled = dbBudget.AlertsEnabled.Bool
	}

	budget := &Budget{
		ID:                dbBudget.ID.String(),
		UserID:            dbBudget.UserID.String(),
		Category:          dbBudget.Category,
		Period:            dbBudget.Period,
		LimitCents:        dbBudget.LimitCents,
		CurrentSpentCents: currentSpentCents,
		AlertThreshold:    alertThreshold,
		AlertsEnabled:     alertsEnabled,
		StartDate:         dbBudget.StartDate,
		EndDate:           dbBudget.EndDate,
	}

	if dbBudget.CreatedAt.Valid {
		budget.CreatedAt = dbBudget.CreatedAt.Time
	}

	return budget
}

// dbBudgetToBudgetSummary converts a database budget to budget summary
func dbBudgetToBudgetSummary(dbBudget *db.Budget) *BudgetSummary {
	currentSpentCents := int64(0)
	if dbBudget.CurrentSpentCents.Valid {
		currentSpentCents = dbBudget.CurrentSpentCents.Int64
	}

	alertThreshold := int16(75) // Default
	if dbBudget.AlertThreshold.Valid {
		alertThreshold = dbBudget.AlertThreshold.Int16
	}

	percentageUsed := 0.0
	if dbBudget.LimitCents > 0 {
		percentageUsed = (float64(currentSpentCents) / float64(dbBudget.LimitCents)) * 100
	}

	isOverBudget := currentSpentCents > dbBudget.LimitCents
	isNearLimit := percentageUsed >= float64(alertThreshold)

	return &BudgetSummary{
		ID:                dbBudget.ID.String(),
		Category:          dbBudget.Category,
		Period:            dbBudget.Period,
		LimitCents:        dbBudget.LimitCents,
		CurrentSpentCents: currentSpentCents,
		PercentageUsed:    percentageUsed,
		AlertThreshold:    alertThreshold,
		IsOverBudget:      isOverBudget,
		IsNearLimit:       isNearLimit,
	}
}

// budgetsToBudgetSummaries converts multiple database budgets to budget summaries
func budgetsToBudgetSummaries(dbBudgets []db.Budget) []*BudgetSummary {
	summaries := make([]*BudgetSummary, len(dbBudgets))
	for i, dbBudget := range dbBudgets {
		summaries[i] = dbBudgetToBudgetSummary(&dbBudget)
	}
	return summaries
}
