package bills

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

// Repository handles data access for bills
type Repository struct {
	db      *sql.DB
	queries *db.Queries
}

// NewRepository creates a new bill repository
func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db:      database,
		queries: db.New(database),
	}
}

// Create creates a new bill
func (r *Repository) Create(ctx context.Context, params db.CreateBillParams) (*db.Bill, error) {
	bill, err := r.queries.CreateBill(ctx, params)
	if err != nil {
		return nil, err
	}
	return &bill, nil
}

// GetByID retrieves a bill by ID
func (r *Repository) GetByID(ctx context.Context, id string) (*db.Bill, error) {
	billID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	bill, err := r.queries.GetBillByID(ctx, billID)
	if err != nil {
		return nil, err
	}

	return &bill, nil
}

// GetByBarcode retrieves a bill by barcode
func (r *Repository) GetByBarcode(ctx context.Context, barcode string) (*db.Bill, error) {
	bill, err := r.queries.GetBillByBarcode(ctx, barcode)
	if err != nil {
		return nil, err
	}

	return &bill, nil
}

// List retrieves bills for a user with pagination
func (r *Repository) List(ctx context.Context, userID string, limit, offset int32) ([]db.Bill, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	bills, err := r.queries.ListUserBills(ctx, db.ListUserBillsParams{
		UserID: userUUID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return bills, nil
}

// ListByStatus retrieves bills for a user filtered by status
func (r *Repository) ListByStatus(ctx context.Context, userID, status string, limit, offset int32) ([]db.Bill, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	bills, err := r.queries.ListUserBillsByStatus(ctx, db.ListUserBillsByStatusParams{
		UserID: userUUID,
		Status: status,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return bills, nil
}

// Count counts total bills for a user
func (r *Repository) Count(ctx context.Context, userID string) (int64, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return 0, err
	}

	count, err := r.queries.CountUserBills(ctx, userUUID)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// UpdateStatus updates a bill's status
func (r *Repository) UpdateStatus(ctx context.Context, id, status string) (*db.Bill, error) {
	billID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	bill, err := r.queries.UpdateBillStatus(ctx, db.UpdateBillStatusParams{
		ID:     billID,
		Status: status,
	})
	if err != nil {
		return nil, err
	}

	return &bill, nil
}

// MarkAsPaid marks a bill as paid
func (r *Repository) MarkAsPaid(ctx context.Context, id string) (*db.Bill, error) {
	billID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	bill, err := r.queries.MarkBillAsPaid(ctx, billID)
	if err != nil {
		return nil, err
	}

	return &bill, nil
}

// Delete soft deletes a bill (marks as cancelled)
func (r *Repository) Delete(ctx context.Context, id string) error {
	billID, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteBill(ctx, billID)
}
