package cards

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/lauratech/fin/back/internal/shared/crypto"
	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

// Repository handles card data access with encryption/decryption
type Repository struct {
	queries       *db.Queries
	encryptionKey []byte
}

// NewRepository creates a new card repository with encryption key
func NewRepository(database *sql.DB, encryptionKey string) *Repository {
	return &Repository{
		queries:       db.New(database),
		encryptionKey: []byte(encryptionKey),
	}
}

// Create creates a new card (encrypts sensitive data before storage)
func (r *Repository) Create(ctx context.Context, params CreateCardParams) (*Card, error) {
	// 1. Encrypt card number
	cardNumberEncrypted, err := crypto.EncryptString(params.CardNumber, r.encryptionKey)
	if err != nil {
		return nil, ErrEncryptionFailed
	}

	// 2. Encrypt CVV
	cvvEncrypted, err := crypto.EncryptString(params.CVV, r.encryptionKey)
	if err != nil {
		return nil, ErrEncryptionFailed
	}

	// 3. Hash PIN if provided
	var pinHash sql.NullString
	if params.PIN != "" {
		hash, err := crypto.HashPIN(params.PIN)
		if err != nil {
			return nil, err
		}
		pinHash = sql.NullString{String: hash, Valid: true}
	}

	// 4. Extract last 4 digits for unencrypted storage (for display)
	lastFour := params.CardNumber[len(params.CardNumber)-4:]

	// 5. Call SQLC generated query
	dbCard, err := r.queries.CreateCard(ctx, db.CreateCardParams{
		UserID:                   uuid.MustParse(params.UserID),
		Type:                     params.Type,
		Brand:                    params.Brand,
		Status:                   "active", // Default status
		CardNumberEncrypted:      cardNumberEncrypted,
		CvvEncrypted:             cvvEncrypted,
		PinHash:                  pinHash,
		LastFourDigits:           lastFour,
		HolderName:               params.HolderName,
		ExpiryMonth:              int16(params.ExpiryMonth),
		ExpiryYear:               int16(params.ExpiryYear),
		DailyLimitCents:          sql.NullInt64{Int64: params.DailyLimitCents, Valid: true},
		MonthlyLimitCents:        sql.NullInt64{Int64: params.MonthlyLimitCents, Valid: true},
		CurrentDailySpentCents:   sql.NullInt64{Int64: 0, Valid: true},
		CurrentMonthlySpentCents: sql.NullInt64{Int64: 0, Valid: true},
		IsContactless:            sql.NullBool{Bool: params.IsContactless, Valid: true},
		IsInternational:          sql.NullBool{Bool: params.IsInternational, Valid: true},
		BlockInternational:       sql.NullBool{Bool: params.BlockInternational, Valid: true},
		BlockOnline:              sql.NullBool{Bool: params.BlockOnline, Valid: true},
		ExpiresAt:                sql.NullTime{Time: params.ExpiresAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// 6. Convert to domain model (with decrypted data for return)
	return dbCardToCard(&dbCard, params.CardNumber, params.CVV), nil
}

// GetByID retrieves a card by ID (decrypts sensitive data)
func (r *Repository) GetByID(ctx context.Context, cardID string) (*Card, error) {
	// 1. Fetch from database
	dbCard, err := r.queries.GetCardByID(ctx, uuid.MustParse(cardID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCardNotFound
		}
		return nil, err
	}

	// 2. Decrypt card number
	cardNumber, err := crypto.DecryptString(dbCard.CardNumberEncrypted, r.encryptionKey)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	// 3. Decrypt CVV
	cvv, err := crypto.DecryptString(dbCard.CvvEncrypted, r.encryptionKey)
	if err != nil {
		return nil, ErrDecryptionFailed
	}

	// 4. Map to domain model (including decrypted data)
	return dbCardToCard(&dbCard, cardNumber, cvv), nil
}

// GetByIDForSummary retrieves a card without decrypting sensitive data
// Use this for list operations where decrypted data is not needed
func (r *Repository) GetByIDForSummary(ctx context.Context, cardID string) (*CardSummary, error) {
	dbCard, err := r.queries.GetCardByID(ctx, uuid.MustParse(cardID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCardNotFound
		}
		return nil, err
	}

	return dbCardToCardSummary(&dbCard), nil
}

// ListUserCards lists all cards for a user (decrypts sensitive data)
func (r *Repository) ListUserCards(ctx context.Context, userID string) ([]Card, error) {
	// 1. Fetch from database
	dbCards, err := r.queries.ListUserCards(ctx, uuid.MustParse(userID))
	if err != nil {
		return nil, err
	}

	// 2. Decrypt each card
	cards := make([]Card, len(dbCards))
	for i, dbCard := range dbCards {
		// Decrypt card number
		cardNumber, err := crypto.DecryptString(dbCard.CardNumberEncrypted, r.encryptionKey)
		if err != nil {
			return nil, ErrDecryptionFailed
		}

		// Decrypt CVV
		cvv, err := crypto.DecryptString(dbCard.CvvEncrypted, r.encryptionKey)
		if err != nil {
			return nil, ErrDecryptionFailed
		}

		cards[i] = *dbCardToCard(&dbCard, cardNumber, cvv)
	}

	return cards, nil
}

// ListUserCardSummaries lists all cards for a user without decrypting sensitive data
// Use this for list responses to avoid unnecessary decryption
func (r *Repository) ListUserCardSummaries(ctx context.Context, userID string) ([]CardSummary, error) {
	dbCards, err := r.queries.ListUserCards(ctx, uuid.MustParse(userID))
	if err != nil {
		return nil, err
	}

	return dbCardsToCardSummaries(dbCards), nil
}

// UpdateStatus updates card status
func (r *Repository) UpdateStatus(ctx context.Context, cardID, status string) error {
	err := r.queries.UpdateCardStatus(ctx, db.UpdateCardStatusParams{
		ID:     uuid.MustParse(cardID),
		Status: status,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrCardNotFound
		}
		return err
	}
	return nil
}

// UpdateLimits updates spending limits
func (r *Repository) UpdateLimits(ctx context.Context, cardID string, params UpdateLimitsParams) error {
	err := r.queries.UpdateCardLimits(ctx, db.UpdateCardLimitsParams{
		ID:                uuid.MustParse(cardID),
		DailyLimitCents:   sql.NullInt64{Int64: params.DailyLimitCents, Valid: true},
		MonthlyLimitCents: sql.NullInt64{Int64: params.MonthlyLimitCents, Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrCardNotFound
		}
		return err
	}
	return nil
}

// UpdateSecuritySettings updates security settings
func (r *Repository) UpdateSecuritySettings(ctx context.Context, cardID string, params SecuritySettingsParams) error {
	err := r.queries.UpdateCardSecuritySettings(ctx, db.UpdateCardSecuritySettingsParams{
		ID:                 uuid.MustParse(cardID),
		IsContactless:      sql.NullBool{Bool: params.IsContactless, Valid: true},
		IsInternational:    sql.NullBool{Bool: params.IsInternational, Valid: true},
		BlockInternational: sql.NullBool{Bool: params.BlockInternational, Valid: true},
		BlockOnline:        sql.NullBool{Bool: params.BlockOnline, Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrCardNotFound
		}
		return err
	}
	return nil
}

// UpdatePIN updates card PIN (hashes before storing)
func (r *Repository) UpdatePIN(ctx context.Context, cardID, pin string) error {
	// Hash the PIN
	hash, err := crypto.HashPIN(pin)
	if err != nil {
		return err
	}

	err = r.queries.UpdateCardPIN(ctx, db.UpdateCardPINParams{
		ID:      uuid.MustParse(cardID),
		PinHash: sql.NullString{String: hash, Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrCardNotFound
		}
		return err
	}
	return nil
}

// VerifyPIN verifies a PIN against the stored hash
func (r *Repository) VerifyPIN(ctx context.Context, cardID, pin string) (bool, error) {
	// Fetch card to get PIN hash
	dbCard, err := r.queries.GetCardByID(ctx, uuid.MustParse(cardID))
	if err != nil {
		if err == sql.ErrNoRows {
			return false, ErrCardNotFound
		}
		return false, err
	}

	// Check if PIN is set
	if !dbCard.PinHash.Valid || dbCard.PinHash.String == "" {
		return false, ErrPINNotSet
	}

	// Verify PIN
	match, err := crypto.VerifyPIN(pin, dbCard.PinHash.String)
	if err != nil {
		return false, err
	}

	return match, nil
}

// GetForUpdate retrieves a card with pessimistic lock (for transactions)
func (r *Repository) GetForUpdate(ctx context.Context, tx *sql.Tx, cardID string) (*db.Card, error) {
	// Create queries from transaction
	qtx := r.queries.WithTx(tx)

	dbCard, err := qtx.GetCardForUpdate(ctx, uuid.MustParse(cardID))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrCardNotFound
		}
		return nil, err
	}

	return &dbCard, nil
}

// UpdateSpentAmounts updates spent amounts (for transaction processing)
func (r *Repository) UpdateSpentAmounts(ctx context.Context, cardID string, dailySpent, monthlySpent int64) error {
	err := r.queries.UpdateCardSpentAmounts(ctx, db.UpdateCardSpentAmountsParams{
		ID:                       uuid.MustParse(cardID),
		CurrentDailySpentCents:   sql.NullInt64{Int64: dailySpent, Valid: true},
		CurrentMonthlySpentCents: sql.NullInt64{Int64: monthlySpent, Valid: true},
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrCardNotFound
		}
		return err
	}
	return nil
}

// CancelCard cancels a card (soft delete)
func (r *Repository) CancelCard(ctx context.Context, cardID string) error {
	err := r.queries.DeleteCard(ctx, uuid.MustParse(cardID))
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrCardNotFound
		}
		return err
	}
	return nil
}

// CountUserCards counts total cards for a user
func (r *Repository) CountUserCards(ctx context.Context, userID string) (int64, error) {
	count, err := r.queries.CountUserCards(ctx, uuid.MustParse(userID))
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CountUserActiveCards counts active cards for a user
func (r *Repository) CountUserActiveCards(ctx context.Context, userID string) (int64, error) {
	count, err := r.queries.CountUserActiveCards(ctx, uuid.MustParse(userID))
	if err != nil {
		return 0, err
	}
	return count, nil
}

// CreateCardTransaction persists a card transaction
func (r *Repository) CreateCardTransaction(
	ctx context.Context,
	params db.CreateCardTransactionParams,
) (*db.CardTransaction, error) {
	tx, err := r.queries.CreateCardTransaction(ctx, params)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// ListCardTransactions retrieves transactions for a card with pagination
func (r *Repository) ListCardTransactions(
	ctx context.Context,
	cardID uuid.UUID,
	limit, offset int32,
) ([]db.CardTransaction, error) {
	txs, err := r.queries.ListCardTransactions(ctx, db.ListCardTransactionsParams{
		CardID: cardID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return txs, nil
}

// ListUserCardTransactions retrieves all user card transactions
func (r *Repository) ListUserCardTransactions(
	ctx context.Context,
	userID uuid.UUID,
	limit, offset int32,
) ([]db.CardTransaction, error) {
	txs, err := r.queries.ListUserCardTransactions(ctx, db.ListUserCardTransactionsParams{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}
	return txs, nil
}

// GetCardTransactionsByCategory retrieves spending by category
func (r *Repository) GetCardTransactionsByCategory(
	ctx context.Context,
	userID uuid.UUID,
	startDate, endDate sql.NullTime,
) ([]db.GetCardTransactionsByCategoryRow, error) {
	return r.queries.GetCardTransactionsByCategory(ctx, db.GetCardTransactionsByCategoryParams{
		UserID:    userID,
		StartDate: startDate.Time,
		EndDate:   endDate.Time,
	})
}

// GetCardTransactionsByDateRange retrieves transactions in a date range
func (r *Repository) GetCardTransactionsByDateRange(
	ctx context.Context,
	userID uuid.UUID,
	startDate, endDate sql.NullTime,
) ([]db.CardTransaction, error) {
	return r.queries.GetCardTransactionsByDateRange(ctx, db.GetCardTransactionsByDateRangeParams{
		UserID:    userID,
		StartDate: startDate.Time,
		EndDate:   endDate.Time,
	})
}
