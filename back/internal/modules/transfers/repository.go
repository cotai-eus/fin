package transfers

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

// Repository handles data access for transfers
type Repository struct {
	db      *sql.DB
	queries *db.Queries
}

// NewRepository creates a new transfer repository
func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db:      database,
		queries: db.New(database),
	}
}

// Create creates a new transfer
func (r *Repository) Create(ctx context.Context, params db.CreateTransferParams) (*db.Transfer, error) {
	transfer, err := r.queries.CreateTransfer(ctx, params)
	if err != nil {
		return nil, err
	}
	return &transfer, nil
}

// GetByID retrieves a transfer by ID
func (r *Repository) GetByID(ctx context.Context, id string) (*db.Transfer, error) {
	transferID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	transfer, err := r.queries.GetTransferByID(ctx, transferID)
	if err != nil {
		return nil, err
	}

	return &transfer, nil
}

// List retrieves transfers for a user with pagination
func (r *Repository) List(ctx context.Context, userID string, limit, offset int32) ([]db.Transfer, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	transfers, err := r.queries.ListUserTransfers(ctx, db.ListUserTransfersParams{
		UserID: userUUID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return transfers, nil
}

// Count counts total transfers for a user
func (r *Repository) Count(ctx context.Context, userID string) (int64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}

	count, err := r.queries.CountUserTransfers(ctx, userUUID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// UpdateStatus updates a transfer's status
func (r *Repository) UpdateStatus(ctx context.Context, id string, status string, failureReason *string) (*db.Transfer, error) {
	transferID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	var nullFailureReason sql.NullString
	if failureReason != nil {
		nullFailureReason = sql.NullString{String: *failureReason, Valid: true}
	}

	transfer, err := r.queries.UpdateTransferStatus(ctx, db.UpdateTransferStatusParams{
		ID:            transferID,
		Status:        status,
		FailureReason: nullFailureReason,
	})
	if err != nil {
		return nil, err
	}

	return &transfer, nil
}

// Cancel cancels a pending transfer
func (r *Repository) Cancel(ctx context.Context, id string) (*db.Transfer, error) {
	transferID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	transfer, err := r.queries.CancelTransfer(ctx, transferID)
	if err != nil {
		return nil, err
	}

	return &transfer, nil
}

// GetDailySum returns the sum of transfers made today by a user
func (r *Repository) GetDailySum(ctx context.Context, userID string) (int64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}

	return r.queries.GetDailyTransferSum(ctx, userUUID)
}

// GetMonthlySum returns the sum of transfers made this month by a user
func (r *Repository) GetMonthlySum(ctx context.Context, userID string) (int64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}

	return r.queries.GetMonthlyTransferSum(ctx, userUUID)
}

// GetForUpdate retrieves a transfer with a pessimistic lock (FOR UPDATE)
func (r *Repository) GetForUpdate(ctx context.Context, id string) (*db.Transfer, error) {
	transferID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	transfer, err := r.queries.GetTransferForUpdate(ctx, transferID)
	if err != nil {
		return nil, err
	}

	return &transfer, nil
}

// WithTx returns a new queries instance bound to a transaction
func (r *Repository) WithTx(tx *sql.Tx) *db.Queries {
	return r.queries.WithTx(tx)
}
