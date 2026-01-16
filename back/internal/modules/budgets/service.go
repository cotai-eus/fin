package budgets

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

// Service handles business logic for budgets
type Service struct {
	repo *Repository
	db   *sql.DB
}

// NewService creates a new budget service
func NewService(repo *Repository, database *sql.DB) *Service {
	return &Service{
		repo: repo,
		db:   database,
	}
}

// CreateBudget creates a new budget
func (s *Service) CreateBudget(ctx context.Context, userID string, req CreateBudgetRequest) (*Budget, error) {
	// Validate category
	if err := ValidateCategory(req.Category); err != nil {
		return nil, err
	}

	// Validate period
	if err := ValidatePeriod(req.Period); err != nil {
		return nil, err
	}

	// Validate limit
	if err := ValidateLimit(req.LimitCents); err != nil {
		return nil, err
	}

	// Set default alert threshold
	alertThreshold := req.AlertThreshold
	if alertThreshold == 0 {
		alertThreshold = 75 // Default: 75%
	}

	// Validate threshold
	if err := ValidateThreshold(alertThreshold); err != nil {
		return nil, err
	}

	// Calculate date range for period
	startDate, endDate, err := CalculateDateRange(req.Period)
	if err != nil {
		return nil, err
	}

	// Check if budget already exists for this category/period
	_, err = s.repo.GetByCategoryAndPeriod(ctx, userID, req.Category, req.Period)
	if err == nil {
		return nil, ErrBudgetAlreadyExists
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Create budget
	userUUID, _ := uuid.Parse(userID)
	dbBudget, err := s.repo.Create(ctx, db.CreateBudgetParams{
		UserID:            userUUID,
		Category:          req.Category,
		Period:            req.Period,
		LimitCents:        req.LimitCents,
		CurrentSpentCents: sql.NullInt64{Int64: 0, Valid: true},
		AlertThreshold:    sql.NullInt16{Int16: alertThreshold, Valid: true},
		AlertsEnabled:     sql.NullBool{Bool: req.AlertsEnabled, Valid: true},
		StartDate:         startDate,
		EndDate:           endDate,
	})
	if err != nil {
		return nil, err
	}

	return dbBudgetToBudget(dbBudget), nil
}

// GetBudgetByID retrieves a budget by ID
func (s *Service) GetBudgetByID(ctx context.Context, userID, budgetID string) (*Budget, error) {
	dbBudget, err := s.repo.GetByID(ctx, budgetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBudgetNotFound
		}
		return nil, err
	}

	// Verify ownership
	if dbBudget.UserID.String() != userID {
		return nil, ErrUnauthorized
	}

	return dbBudgetToBudget(dbBudget), nil
}

// ListUserBudgets retrieves budgets for a user with pagination
func (s *Service) ListUserBudgets(ctx context.Context, userID string, params BudgetListParams) ([]*BudgetSummary, int64, error) {
	// Set defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 20
	}

	offset := int32((params.Page - 1) * params.Limit)
	limit := int32(params.Limit)

	var dbBudgets []db.Budget
	var err error

	// List budgets (with optional category filter)
	if params.Category != "" {
		dbBudgets, err = s.repo.ListByCategory(ctx, userID, params.Category, limit, offset)
	} else {
		dbBudgets, err = s.repo.List(ctx, userID, limit, offset)
	}

	if err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := s.repo.Count(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return budgetsToBudgetSummaries(dbBudgets), total, nil
}

// UpdateBudget updates a budget
func (s *Service) UpdateBudget(ctx context.Context, userID, budgetID string, req UpdateBudgetRequest) (*Budget, error) {
	// Get existing budget
	dbBudget, err := s.repo.GetByID(ctx, budgetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBudgetNotFound
		}
		return nil, err
	}

	// Verify ownership
	if dbBudget.UserID.String() != userID {
		return nil, ErrUnauthorized
	}

	// Prepare update params
	budgetUUID, _ := uuid.Parse(budgetID)

	limitCents := dbBudget.LimitCents
	if req.LimitCents > 0 {
		if err := ValidateLimit(req.LimitCents); err != nil {
			return nil, err
		}
		limitCents = req.LimitCents
	}

	alertThreshold := int16(75)
	if dbBudget.AlertThreshold.Valid {
		alertThreshold = dbBudget.AlertThreshold.Int16
	}
	if req.AlertThreshold > 0 {
		if err := ValidateThreshold(req.AlertThreshold); err != nil {
			return nil, err
		}
		alertThreshold = req.AlertThreshold
	}

	alertsEnabled := true
	if dbBudget.AlertsEnabled.Valid {
		alertsEnabled = dbBudget.AlertsEnabled.Bool
	}
	if req.AlertsEnabled != nil {
		alertsEnabled = *req.AlertsEnabled
	}

	updateParams := db.UpdateBudgetParams{
		ID:             budgetUUID,
		LimitCents:     limitCents,
		AlertThreshold: sql.NullInt16{Int16: alertThreshold, Valid: true},
		AlertsEnabled:  sql.NullBool{Bool: alertsEnabled, Valid: true},
	}

	// Update budget
	updatedBudget, err := s.repo.Update(ctx, updateParams)
	if err != nil {
		return nil, err
	}

	return dbBudgetToBudget(updatedBudget), nil
}

// DeleteBudget deletes a budget
func (s *Service) DeleteBudget(ctx context.Context, userID, budgetID string) error {
	// Get budget
	dbBudget, err := s.repo.GetByID(ctx, budgetID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrBudgetNotFound
		}
		return err
	}

	// Verify ownership
	if dbBudget.UserID.String() != userID {
		return ErrUnauthorized
	}

	// Delete budget
	return s.repo.Delete(ctx, budgetID)
}

// GetBudgetSummary retrieves a summary of all budgets for a user
func (s *Service) GetBudgetSummary(ctx context.Context, userID string) (*BudgetAnalytics, error) {
	// Get all budgets
	dbBudgets, err := s.repo.List(ctx, userID, 100, 0) // Get up to 100 budgets
	if err != nil {
		return nil, err
	}

	totalBudgetCents := int64(0)
	totalSpentCents := int64(0)
	categoriesOverBudget := 0
	categoryBreakdown := make([]*CategorySpending, 0)

	for _, dbBudget := range dbBudgets {
		currentSpentCents := int64(0)
		if dbBudget.CurrentSpentCents.Valid {
			currentSpentCents = dbBudget.CurrentSpentCents.Int64
		}

		totalBudgetCents += dbBudget.LimitCents
		totalSpentCents += currentSpentCents

		if currentSpentCents > dbBudget.LimitCents {
			categoriesOverBudget++
		}

		percentage := 0.0
		if dbBudget.LimitCents > 0 {
			percentage = (float64(currentSpentCents) / float64(dbBudget.LimitCents)) * 100
		}

		categoryBreakdown = append(categoryBreakdown, &CategorySpending{
			Category:    dbBudget.Category,
			TotalCents:  currentSpentCents,
			Percentage:  percentage,
			BudgetCents: dbBudget.LimitCents,
		})
	}

	percentageUsed := 0.0
	if totalBudgetCents > 0 {
		percentageUsed = (float64(totalSpentCents) / float64(totalBudgetCents)) * 100
	}

	return &BudgetAnalytics{
		TotalBudgetCents:     totalBudgetCents,
		TotalSpentCents:      totalSpentCents,
		PercentageUsed:       percentageUsed,
		CategoriesOverBudget: categoriesOverBudget,
		CategoryBreakdown:    categoryBreakdown,
	}, nil
}
