package transfers

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/lauratech/fin/back/internal/modules/users"
	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

const (
	TEDFeeCents = 1000 // R$ 10.00 TED fee
)

// Service handles business logic for transfers
type Service struct {
	repo     *Repository
	userRepo *users.Repository
	db       *sql.DB
}

// NewService creates a new transfer service
func NewService(repo *Repository, userRepo *users.Repository, database *sql.DB) *Service {
	return &Service{
		repo:     repo,
		userRepo: userRepo,
		db:       database,
	}
}

// ExecutePIX executes a PIX transfer with balance and limit validation
func (s *Service) ExecutePIX(ctx context.Context, userID string, req CreatePIXRequest) (*Transfer, error) {
	// Validate PIX key
	if err := ValidatePIXKey(req.PixKey, req.PixKeyType); err != nil {
		return nil, err
	}

	// Validate amount
	if err := ValidateAmount(req.AmountCents); err != nil {
		return nil, err
	}

	// Execute transfer in transaction
	var transfer *db.Transfer
	err := s.executeInTransaction(ctx, func(tx *sql.Tx) error {
		qtx := db.New(tx)
		userUUID, _ := uuid.Parse(userID)

		// 1. Lock user record (FOR UPDATE)
		user, err := qtx.GetUserForUpdate(ctx, userUUID)
		if err != nil {
			return err
		}

		// 2. Check balance
		userBalance := int64(0)
		if user.BalanceCents.Valid {
			userBalance = user.BalanceCents.Int64
		}
		if userBalance < req.AmountCents {
			return ErrInsufficientBalance
		}

		// 3. Check daily limit
		dailySum, err := s.getDailySumWithTx(ctx, tx, userID)
		if err != nil {
			return err
		}
		dailyLimit := int64(100000) // Default
		if user.DailyTransferLimitCents.Valid {
			dailyLimit = user.DailyTransferLimitCents.Int64
		}
		if dailySum+req.AmountCents > dailyLimit {
			return ErrDailyLimitExceeded
		}

		// 4. Check monthly limit
		monthlySum, err := s.getMonthlySumWithTx(ctx, tx, userID)
		if err != nil {
			return err
		}
		monthlyLimit := int64(500000) // Default
		if user.MonthlyTransferLimitCents.Valid {
			monthlyLimit = user.MonthlyTransferLimitCents.Int64
		}
		if monthlySum+req.AmountCents > monthlyLimit {
			return ErrMonthlyLimitExceeded
		}

		// 5. Debit user balance
		err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
			ID:           userUUID,
			BalanceCents: sql.NullInt64{Int64: -req.AmountCents, Valid: true},
		})
		if err != nil {
			return err
		}

		// 6. Create transfer record
		dbTransfer, err := qtx.CreateTransfer(ctx, db.CreateTransferParams{
			UserID:      userUUID,
			Type:        "pix",
			Status:      "completed",
			AmountCents: req.AmountCents,
			FeeCents:    sql.NullInt64{Int64: 0, Valid: true},
			Currency:    sql.NullString{String: "BRL", Valid: true},
			PixKey:      sql.NullString{String: req.PixKey, Valid: true},
			PixKeyType:  sql.NullString{String: req.PixKeyType, Valid: true},
			CompletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}

		transfer = &dbTransfer
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dbTransferToTransfer(transfer), nil
}

// ExecuteTED executes a TED transfer with R$ 10.00 fee
func (s *Service) ExecuteTED(ctx context.Context, userID string, req CreateTEDRequest) (*Transfer, error) {
	// Validate TED data
	if err := ValidateTEDData(req); err != nil {
		return nil, err
	}

	totalAmount := req.AmountCents + TEDFeeCents

	// Execute transfer in transaction
	var transfer *db.Transfer
	err := s.executeInTransaction(ctx, func(tx *sql.Tx) error {
		qtx := db.New(tx)
		userUUID, _ := uuid.Parse(userID)

		// 1. Lock user record
		user, err := qtx.GetUserForUpdate(ctx, userUUID)
		if err != nil {
			return err
		}

		// 2. Check balance (amount + fee)
		userBalance := int64(0)
		if user.BalanceCents.Valid {
			userBalance = user.BalanceCents.Int64
		}
		if userBalance < totalAmount {
			return ErrInsufficientBalance
		}

		// 3. Check daily limit
		dailySum, err := s.getDailySumWithTx(ctx, tx, userID)
		if err != nil {
			return err
		}
		dailyLimit := int64(100000)
		if user.DailyTransferLimitCents.Valid {
			dailyLimit = user.DailyTransferLimitCents.Int64
		}
		if dailySum+totalAmount > dailyLimit {
			return ErrDailyLimitExceeded
		}

		// 4. Check monthly limit
		monthlySum, err := s.getMonthlySumWithTx(ctx, tx, userID)
		if err != nil {
			return err
		}
		monthlyLimit := int64(500000)
		if user.MonthlyTransferLimitCents.Valid {
			monthlyLimit = user.MonthlyTransferLimitCents.Int64
		}
		if monthlySum+totalAmount > monthlyLimit {
			return ErrMonthlyLimitExceeded
		}

		// 5. Debit user balance
		err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
			ID:           userUUID,
			BalanceCents: sql.NullInt64{Int64: -totalAmount, Valid: true},
		})
		if err != nil {
			return err
		}

		// 6. Create transfer record
		dbTransfer, err := qtx.CreateTransfer(ctx, db.CreateTransferParams{
			UserID:               userUUID,
			Type:                 "ted",
			Status:               "completed",
			AmountCents:          req.AmountCents,
			FeeCents:             sql.NullInt64{Int64: TEDFeeCents, Valid: true},
			Currency:             sql.NullString{String: "BRL", Valid: true},
			RecipientName:        sql.NullString{String: req.RecipientName, Valid: true},
			RecipientDocument:    sql.NullString{String: req.RecipientDocument, Valid: true},
			RecipientBank:        sql.NullString{String: req.RecipientBank, Valid: true},
			RecipientBranch:      sql.NullString{String: req.RecipientBranch, Valid: true},
			RecipientAccount:     sql.NullString{String: req.RecipientAccount, Valid: true},
			RecipientAccountType: sql.NullString{String: req.RecipientAccountType, Valid: true},
			CompletedAt:          sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}

		transfer = &dbTransfer
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dbTransferToTransfer(transfer), nil
}

// ExecuteP2P executes a peer-to-peer transfer between two users
func (s *Service) ExecuteP2P(ctx context.Context, senderID string, req CreateP2PRequest) (*Transfer, error) {
	// Validate amount
	if err := ValidateAmount(req.AmountCents); err != nil {
		return nil, err
	}

	// Validate recipient exists
	_, err := s.userRepo.GetByID(ctx, req.RecipientUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecipientNotFound
		}
		return nil, err
	}

	// Cannot transfer to yourself
	if senderID == req.RecipientUserID {
		return nil, ErrCannotTransferToSelf
	}

	// Execute P2P transfer in transaction
	var transfer *db.Transfer
	err = s.executeInTransaction(ctx, func(tx *sql.Tx) error {
		qtx := db.New(tx)
		senderUUID, _ := uuid.Parse(senderID)
		recipientUUID, _ := uuid.Parse(req.RecipientUserID)

		// 1. Lock sender (order by UUID to prevent deadlocks)
		sender, err := qtx.GetUserForUpdate(ctx, senderUUID)
		if err != nil {
			return err
		}

		// 2. Lock recipient
		_, err = qtx.GetUserForUpdate(ctx, recipientUUID)
		if err != nil {
			return err
		}

		// 3. Check sender balance
		senderBalance := int64(0)
		if sender.BalanceCents.Valid {
			senderBalance = sender.BalanceCents.Int64
		}
		if senderBalance < req.AmountCents {
			return ErrInsufficientBalance
		}

		// 4. Check daily limit
		dailySum, err := s.getDailySumWithTx(ctx, tx, senderID)
		if err != nil {
			return err
		}
		dailyLimit := int64(100000)
		if sender.DailyTransferLimitCents.Valid {
			dailyLimit = sender.DailyTransferLimitCents.Int64
		}
		if dailySum+req.AmountCents > dailyLimit {
			return ErrDailyLimitExceeded
		}

		// 5. Check monthly limit
		monthlySum, err := s.getMonthlySumWithTx(ctx, tx, senderID)
		if err != nil {
			return err
		}
		monthlyLimit := int64(500000)
		if sender.MonthlyTransferLimitCents.Valid {
			monthlyLimit = sender.MonthlyTransferLimitCents.Int64
		}
		if monthlySum+req.AmountCents > monthlyLimit {
			return ErrMonthlyLimitExceeded
		}

		// 6. Debit sender
		err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
			ID:           senderUUID,
			BalanceCents: sql.NullInt64{Int64: -req.AmountCents, Valid: true},
		})
		if err != nil {
			return err
		}

		// 7. Credit recipient
		err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
			ID:           recipientUUID,
			BalanceCents: sql.NullInt64{Int64: req.AmountCents, Valid: true},
		})
		if err != nil {
			return err
		}

		// 8. Create transfer record
		dbTransfer, err := qtx.CreateTransfer(ctx, db.CreateTransferParams{
			UserID:          senderUUID,
			Type:            "p2p",
			Status:          "completed",
			AmountCents:     req.AmountCents,
			FeeCents:        sql.NullInt64{Int64: 0, Valid: true},
			Currency:        sql.NullString{String: "BRL", Valid: true},
			RecipientUserID: uuid.NullUUID{UUID: recipientUUID, Valid: true},
			CompletedAt:     sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}

		transfer = &dbTransfer
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dbTransferToTransfer(transfer), nil
}

// List retrieves transfers for a user with pagination
func (s *Service) List(ctx context.Context, userID string, params TransferListParams) ([]Transfer, int, error) {
	// Calculate offset
	offset := (params.Page - 1) * params.Limit

	// Get transfers
	dbTransfers, err := s.repo.List(ctx, userID, int32(params.Limit), int32(offset))
	if err != nil {
		return nil, 0, err
	}

	// Get total count
	totalCount, err := s.repo.Count(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	// Convert to domain models
	transfers := dbTransfersToTransfers(dbTransfers)

	return transfers, int(totalCount), nil
}

// Cancel cancels a pending transfer and refunds the balance
func (s *Service) Cancel(ctx context.Context, userID, transferID string) error {
	// Execute cancellation in transaction
	return s.executeInTransaction(ctx, func(tx *sql.Tx) error {
		qtx := db.New(tx)
		transferUUID, _ := uuid.Parse(transferID)
		userUUID, _ := uuid.Parse(userID)

		// 1. Get transfer (with lock)
		transfer, err := qtx.GetTransferForUpdate(ctx, transferUUID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return ErrTransferNotFound
			}
			return err
		}

		// 2. Verify ownership
		if transfer.UserID != userUUID {
			return ErrTransferNotFound
		}

		// 3. Verify status (can only cancel pending transfers)
		if transfer.Status != "pending" {
			return ErrInvalidTransferStatus
		}

		// 4. Refund balance (amount + fee)
		refundAmount := transfer.AmountCents
		if transfer.FeeCents.Valid {
			refundAmount += transfer.FeeCents.Int64
		}

		err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
			ID:           userUUID,
			BalanceCents: sql.NullInt64{Int64: refundAmount, Valid: true},
		})
		if err != nil {
			return err
		}

		// 5. Cancel transfer
		_, err = qtx.CancelTransfer(ctx, transferUUID)
		if err != nil {
			return err
		}

		return nil
	})
}

// GetByID retrieves a transfer by ID
func (s *Service) GetByID(ctx context.Context, userID, transferID string) (*Transfer, error) {
	dbTransfer, err := s.repo.GetByID(ctx, transferID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTransferNotFound
		}
		return nil, err
	}

	// Verify ownership
	if dbTransfer.UserID.String() != userID {
		return nil, ErrTransferNotFound
	}

	return dbTransferToTransfer(dbTransfer), nil
}

// ExecuteDeposit processes a deposit transaction
func (s *Service) ExecuteDeposit(ctx context.Context, userID string, req ExecuteDepositRequest) (*Transfer, error) {
	// Validate amount
	if err := ValidateAmount(req.AmountCents); err != nil {
		return nil, err
	}

	// Execute deposit in transaction
	var transfer *db.Transfer
	err := s.executeInTransaction(ctx, func(tx *sql.Tx) error {
		qtx := db.New(tx)
		userUUID, _ := uuid.Parse(userID)

		// 1. Create transfer record
		dbTransfer, err := qtx.CreateTransfer(ctx, db.CreateTransferParams{
			UserID:      userUUID,
			Type:        "deposit",
			Status:      "completed",
			AmountCents: req.AmountCents,
			FeeCents:    sql.NullInt64{Int64: 0, Valid: true},
			Currency:    sql.NullString{String: "BRL", Valid: true},
			CompletedAt: sql.NullTime{Time: time.Now(), Valid: true},
		})
		if err != nil {
			return err
		}

		// 2. Credit user balance
		err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
			ID:           userUUID,
			BalanceCents: sql.NullInt64{Int64: req.AmountCents, Valid: true},
		})
		if err != nil {
			return err
		}

		transfer = &dbTransfer
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dbTransferToTransfer(transfer), nil
}

// CreatePaymentRequest creates a P2P payment request (pending status)
func (s *Service) CreatePaymentRequest(ctx context.Context, userID string, req CreatePaymentRequestRequest) (*Transfer, error) {
	// Validate amount
	if err := ValidateAmount(req.AmountCents); err != nil {
		return nil, err
	}

	// Cannot request payment from self
	if req.RecipientUserID == userID {
		return nil, ErrCannotTransferToSelf
	}

	// Check recipient exists
	recipientUUID, err := uuid.Parse(req.RecipientUserID)
	if err != nil {
		return nil, ErrRecipientNotFound
	}

	_, err = s.userRepo.GetByID(context.Background(), req.RecipientUserID)
	if err != nil {
		return nil, ErrRecipientNotFound
	}

	// Create pending transfer
	userUUID, _ := uuid.Parse(userID)
	dbTransfer, err := s.repo.Create(ctx, db.CreateTransferParams{
		UserID:          userUUID,
		Type:            "p2p",
		Status:          "pending",
		AmountCents:     req.AmountCents,
		FeeCents:        sql.NullInt64{Int64: 0, Valid: true},
		Currency:        sql.NullString{String: "BRL", Valid: true},
		RecipientUserID: uuid.NullUUID{UUID: recipientUUID, Valid: true},
	})

	if err != nil {
		return nil, err
	}

	return dbTransferToTransfer(dbTransfer), nil
}

// ListPaymentRequests lists pending payment requests for a user
func (s *Service) ListPaymentRequests(ctx context.Context, userID string) ([]Transfer, error) {
	userUUID, _ := uuid.Parse(userID)

	// Query transfers with type='p2p' and status='pending'
	dbTransfers, err := s.repo.ListByUserAndStatus(ctx, userUUID, "pending")
	if err != nil {
		return nil, err
	}

	return dbTransfersToTransfers(dbTransfers), nil
}

// ApprovePaymentRequest approves a payment request and executes the transfer
func (s *Service) ApprovePaymentRequest(ctx context.Context, userID string, requestID string) (*Transfer, error) {
	requestUUID, err := uuid.Parse(requestID)
	if err != nil {
		return nil, ErrTransferNotFound
	}

	// Execute in transaction
	var transfer *db.Transfer
	err = s.executeInTransaction(ctx, func(tx *sql.Tx) error {
		qtx := db.New(tx)
		userUUID, _ := uuid.Parse(userID)

		// 1. Get payment request
		dbTransfer, err := qtx.GetTransferByID(ctx, requestUUID)
		if err != nil {
			return ErrTransferNotFound
		}

		// 2. Verify recipient is the approver
		if !dbTransfer.RecipientUserID.Valid || dbTransfer.RecipientUserID.UUID.String() != userID {
			return ErrTransferNotFound
		}

		// 3. Check status is pending
		if dbTransfer.Status != "pending" {
			return ErrInvalidTransferStatus
		}

		// 4. Lock payer user record
		payerUser, err := qtx.GetUserForUpdate(ctx, dbTransfer.UserID)
		if err != nil {
			return err
		}

		// 5. Check balance
		payerBalance := int64(0)
		if payerUser.BalanceCents.Valid {
			payerBalance = payerUser.BalanceCents.Int64
		}
		if payerBalance < dbTransfer.AmountCents {
			return ErrInsufficientBalance
		}

		// 6. Debit payer
		err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
			ID:           dbTransfer.UserID,
			BalanceCents: sql.NullInt64{Int64: -dbTransfer.AmountCents, Valid: true},
		})
		if err != nil {
			return err
		}

		// 7. Credit recipient
		err = qtx.UpdateUserBalance(ctx, db.UpdateUserBalanceParams{
			ID:           userUUID,
			BalanceCents: sql.NullInt64{Int64: dbTransfer.AmountCents, Valid: true},
		})
		if err != nil {
			return err
		}

		// 8. Update transfer status
		updatedTransfer, err := qtx.UpdateTransferStatus(ctx, db.UpdateTransferStatusParams{
			ID:            requestUUID,
			Status:        "completed",
			FailureReason: sql.NullString{Valid: false},
		})
		if err != nil {
			return err
		}

		transfer = &updatedTransfer
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dbTransferToTransfer(transfer), nil
}

// RejectPaymentRequest rejects a payment request
func (s *Service) RejectPaymentRequest(ctx context.Context, userID string, requestID string) error {
	// Get payment request
	dbTransfer, err := s.repo.GetByID(ctx, requestID)
	if err != nil {
		return ErrTransferNotFound
	}

	// Verify recipient is the rejector
	if !dbTransfer.RecipientUserID.Valid || dbTransfer.RecipientUserID.UUID.String() != userID {
		return ErrTransferNotFound
	}

	// Check status is pending
	if dbTransfer.Status != "pending" {
		return ErrInvalidTransferStatus
	}

	// Update status to cancelled
	_, err = s.repo.UpdateStatus(ctx, requestID, "cancelled", nil)
	return err
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
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// getDailySumWithTx gets daily transfer sum within a transaction
func (s *Service) getDailySumWithTx(ctx context.Context, tx *sql.Tx, userID string) (int64, error) {
	qtx := db.New(tx)
	userUUID, _ := uuid.Parse(userID)

	return qtx.GetDailyTransferSum(ctx, userUUID)
}

// getMonthlySumWithTx gets monthly transfer sum within a transaction
func (s *Service) getMonthlySumWithTx(ctx context.Context, tx *sql.Tx, userID string) (int64, error) {
	qtx := db.New(tx)
	userUUID, _ := uuid.Parse(userID)

	return qtx.GetMonthlyTransferSum(ctx, userUUID)
}

// interfaceToInt64 safely converts interface{} to int64
func interfaceToInt64(v interface{}) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case float64:
		return int64(math.Round(val))
	case int:
		return int64(val)
	default:
		return 0
	}
}
