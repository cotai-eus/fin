package bills

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

const (
	BillPaymentFeeCents = 200 // R$ 2.00 fee per bill payment
)

// Service handles business logic for bills
type Service struct {
	repo *Repository
	db   *sql.DB
}

// NewService creates a new bill service
func NewService(repo *Repository, database *sql.DB) *Service {
	return &Service{
		repo: repo,
		db:   database,
	}
}

// ValidateBarcodeInfo validates a barcode and returns information
func (s *Service) ValidateBarcodeInfo(ctx context.Context, barcode string) (*ValidateBarcodeResponse, error) {
	// Validate and parse barcode
	barcodeInfo, err := ValidateBarcode(barcode)
	if err != nil {
		return nil, err
	}

	return &ValidateBarcodeResponse{
		Valid:         true,
		RecipientName: barcodeInfo.RecipientName,
		AmountCents:   barcodeInfo.AmountCents,
		DueDate:       barcodeInfo.DueDate.Format("2006-01-02"),
		Type:          barcodeInfo.BillType,
	}, nil
}

// CreateBill creates a new bill from a barcode
func (s *Service) CreateBill(ctx context.Context, userID string, req CreateBillRequest) (*Bill, error) {
	// Validate barcode
	barcodeInfo, err := ValidateBarcode(req.Barcode)
	if err != nil {
		return nil, err
	}

	// Validate type
	if err := ValidateBillType(req.Type); err != nil {
		return nil, err
	}

	// Check if barcode already exists
	_, err = s.repo.GetByBarcode(ctx, barcodeInfo.Barcode)
	if err == nil {
		return nil, ErrBarcodeAlreadyExists
	} else if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Calculate final amount (amount + fee)
	finalAmountCents := barcodeInfo.AmountCents + BillPaymentFeeCents

	// Create bill
	userUUID, _ := uuid.Parse(userID)
	dbBill, err := s.repo.Create(ctx, db.CreateBillParams{
		UserID:           userUUID,
		Type:             req.Type,
		Status:           "pending",
		Barcode:          barcodeInfo.Barcode,
		AmountCents:      barcodeInfo.AmountCents,
		FeeCents:         sql.NullInt64{Int64: BillPaymentFeeCents, Valid: true},
		FinalAmountCents: finalAmountCents,
		RecipientName:    barcodeInfo.RecipientName,
		DueDate:          barcodeInfo.DueDate,
	})
	if err != nil {
		return nil, err
	}

	return dbBillToBill(dbBill), nil
}

// GetBillByID retrieves a bill by ID
func (s *Service) GetBillByID(ctx context.Context, userID, billID string) (*Bill, error) {
	dbBill, err := s.repo.GetByID(ctx, billID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBillNotFound
		}
		return nil, err
	}

	// Verify ownership
	if dbBill.UserID.String() != userID {
		return nil, ErrUnauthorized
	}

	return dbBillToBill(dbBill), nil
}

// ListUserBills retrieves bills for a user with pagination
func (s *Service) ListUserBills(ctx context.Context, userID string, params BillListParams) ([]*BillSummary, int64, error) {
	// Set defaults
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 || params.Limit > 100 {
		params.Limit = 20
	}

	offset := int32((params.Page - 1) * params.Limit)
	limit := int32(params.Limit)

	var dbBills []db.Bill
	var err error

	// List bills (with optional status filter)
	if params.Status != "" {
		dbBills, err = s.repo.ListByStatus(ctx, userID, params.Status, limit, offset)
	} else {
		dbBills, err = s.repo.List(ctx, userID, limit, offset)
	}

	if err != nil {
		return nil, 0, err
	}

	// Get total count
	total, err := s.repo.Count(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	return billsToBillSummaries(dbBills), total, nil
}

// PayBill processes a bill payment
func (s *Service) PayBill(ctx context.Context, userID, billID string) (*Bill, error) {
	var bill *db.Bill

	err := s.executeInTransaction(ctx, func(tx *sql.Tx) error {
		qtx := db.New(tx)
		userUUID, _ := uuid.Parse(userID)
		billUUID, _ := uuid.Parse(billID)

		// 1. Lock and get bill
		dbBill, err := qtx.GetBillForUpdate(ctx, billUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrBillNotFound
			}
			return err
		}

		// 2. Verify ownership
		if dbBill.UserID != userUUID {
			return ErrUnauthorized
		}

		// 3. Check bill status
		if dbBill.Status == "paid" {
			return ErrBillAlreadyPaid
		}
		if dbBill.Status == "cancelled" {
			return ErrBillCancelled
		}

		// 4. Lock user record
		user, err := qtx.GetUserForUpdate(ctx, userUUID)
		if err != nil {
			return err
		}

		// 5. Check balance
		userBalance := int64(0)
		if user.BalanceCents.Valid {
			userBalance = user.BalanceCents.Int64
		}
		if userBalance < dbBill.FinalAmountCents {
			return ErrInsufficientBalance
		}

		// 6. Debit user balance
		err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
			ID:           userUUID,
			BalanceCents: sql.NullInt64{Int64: -dbBill.FinalAmountCents, Valid: true},
		})
		if err != nil {
			return err
		}

		// 7. Mark bill as paid
		paidBill, err := qtx.MarkBillAsPaid(ctx, billUUID)
		if err != nil {
			return err
		}

		bill = &paidBill
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dbBillToBill(bill), nil
}

// CancelBill cancels a bill
func (s *Service) CancelBill(ctx context.Context, userID, billID string) error {
	// Get bill
	dbBill, err := s.repo.GetByID(ctx, billID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrBillNotFound
		}
		return err
	}

	// Verify ownership
	if dbBill.UserID.String() != userID {
		return ErrUnauthorized
	}

	// Check if bill is already paid
	if dbBill.Status == "paid" {
		return ErrBillAlreadyPaid
	}

	// Cancel bill
	return s.repo.Delete(ctx, billID)
}

// executeInTransaction executes a function within a database transaction
func (s *Service) executeInTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	err = fn(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
