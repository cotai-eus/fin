package budgets

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

// Repository handles data access for budgets
type Repository struct {
	db      *sql.DB
	queries *db.Queries
}

// NewRepository creates a new budget repository
func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db:      database,
		queries: db.New(database),
	}
}

// Create creates a new budget
func (r *Repository) Create(ctx context.Context, params db.CreateBudgetParams) (*db.Budget, error) {
	budget, err := r.queries.CreateBudget(ctx, params)
	if err != nil {
		return nil, err
	}
	return &budget, nil
}

// GetByID retrieves a budget by ID
func (r *Repository) GetByID(ctx context.Context, id string) (*db.Budget, error) {
	budgetID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	budget, err := r.queries.GetBudgetByID(ctx, budgetID)
	if err != nil {
		return nil, err
	}

	return &budget, nil
}

// GetByCategoryAndPeriod retrieves a budget by category and period
func (r *Repository) GetByCategoryAndPeriod(ctx context.Context, userID, category, period string) (*db.Budget, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	budget, err := r.queries.GetBudgetByCategoryAndPeriod(ctx, db.GetBudgetByCategoryAndPeriodParams{
		UserID:   userUUID,
		Category: category,
		Period:   period,
	})
	if err != nil {
		return nil, err
	}

	return &budget, nil
}

// List retrieves budgets for a user with pagination
func (r *Repository) List(ctx context.Context, userID string, limit, offset int32) ([]db.Budget, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	budgets, err := r.queries.ListUserBudgets(ctx, db.ListUserBudgetsParams{
		UserID: userUUID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return budgets, nil
}

// ListByCategory retrieves budgets for a user filtered by category
func (r *Repository) ListByCategory(ctx context.Context, userID, category string, limit, offset int32) ([]db.Budget, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	budgets, err := r.queries.ListUserBudgetsByCategory(ctx, db.ListUserBudgetsByCategoryParams{
		UserID:   userUUID,
		Category: category,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		return nil, err
	}

	return budgets, nil
}

// Count counts total budgets for a user
func (r *Repository) Count(ctx context.Context, userID string) (int64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}

	count, err := r.queries.CountUserBudgets(ctx, userUUID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// Update updates a budget
func (r *Repository) Update(ctx context.Context, params db.UpdateBudgetParams) (*db.Budget, error) {
	budget, err := r.queries.UpdateBudget(ctx, params)
	if err != nil {
		return nil, err
	}

	return &budget, nil
}

// UpdateSpent updates the current spent amount
func (r *Repository) UpdateSpent(ctx context.Context, id string, currentSpentCents int64) (*db.Budget, error) {
	budgetID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	budget, err := r.queries.UpdateBudgetSpent(ctx, db.UpdateBudgetSpentParams{
		ID: budgetID,
		CurrentSpentCents: sql.NullInt64{
			Int64: currentSpentCents,
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return &budget, nil
}

// Delete deletes a budget
func (r *Repository) Delete(ctx context.Context, id string) error {
	budgetID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteBudget(ctx, budgetID)
}

// ResetSpent resets all budgets' spent amounts (for period rollover)
func (r *Repository) ResetSpent(ctx context.Context, userID string) error {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	return r.queries.ResetBudgetSpent(ctx, userUUID)
}
