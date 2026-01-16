package cards

import (
	"context"
	"database/sql"
	"time"
)

// Service handles card business logic
type Service struct {
	repo *Repository
	db   *sql.DB
}

// NewService creates a new card service
func NewService(repo *Repository, database *sql.DB) *Service {
	return &Service{
		repo: repo,
		db:   database,
	}
}

// CreateCard creates a new physical or virtual card
func (s *Service) CreateCard(ctx context.Context, userID string, req CreateCardRequest) (*Card, error) {
	// 1. Validate card type
	if err := ValidateCardType(req.Type); err != nil {
		return nil, err
	}

	// 2. Validate card brand
	if err := ValidateCardBrand(req.Brand); err != nil {
		return nil, err
	}

	// 3. Validate or generate card number
	var cardNumber string
	if req.CardNumber != "" {
		// Validate provided card number
		if err := ValidateCardNumber(req.CardNumber); err != nil {
			return nil, err
		}
		cardNumber = req.CardNumber
	} else {
		// Auto-generate card number
		generated, err := GenerateCardNumber(req.Brand)
		if err != nil {
			return nil, err
		}
		cardNumber = generated
	}

	// 4. Validate CVV format
	if err := ValidateCVV(req.CVV); err != nil {
		return nil, err
	}

	// 5. Validate PIN if provided
	if req.PIN != "" {
		if err := ValidatePIN(req.PIN); err != nil {
			return nil, err
		}
	}

	// 6. Set default limits if not provided
	dailyLimit := req.DailyLimitCents
	if dailyLimit == 0 {
		dailyLimit = 500000 // R$ 5,000 default
	}

	monthlyLimit := req.MonthlyLimitCents
	if monthlyLimit == 0 {
		monthlyLimit = 5000000 // R$ 50,000 default
	}

	// 7. Validate limits
	if dailyLimit < 0 || monthlyLimit < 0 {
		return nil, ErrInvalidLimit
	}

	// 8. Calculate expiry date
	expiresAt := CalculateExpiryDate(req.Type)
	expiryMonth := int(expiresAt.Month())
	expiryYear := expiresAt.Year()

	// 9. Call repository (which handles encryption)
	return s.repo.Create(ctx, CreateCardParams{
		UserID:             userID,
		Type:               req.Type,
		Brand:              req.Brand,
		CardNumber:         cardNumber,
		CVV:                req.CVV,
		PIN:                req.PIN,
		HolderName:         req.HolderName,
		ExpiryMonth:        expiryMonth,
		ExpiryYear:         expiryYear,
		DailyLimitCents:    dailyLimit,
		MonthlyLimitCents:  monthlyLimit,
		IsContactless:      true,  // Default enabled
		IsInternational:    false, // Default disabled for security
		BlockInternational: false,
		BlockOnline:        false,
		ExpiresAt:          expiresAt,
	})
}

// GetCardByID retrieves card details (returns masked card number)
func (s *Service) GetCardByID(ctx context.Context, userID, cardID string) (*CardDetails, error) {
	// 1. Get card with decrypted data
	card, err := s.repo.GetByID(ctx, cardID)
	if err != nil {
		return nil, err
	}

	// 2. Verify ownership
	if card.UserID != userID {
		return nil, ErrUnauthorized
	}

	// 3. Convert to card details (masks card number)
	// Note: hasPIN check would require additional repository method
	hasPIN := false // TODO: Implement proper PIN check
	return cardToCardDetails(card, hasPIN), nil
}

// ListUserCards lists all cards for a user (without sensitive data)
func (s *Service) ListUserCards(ctx context.Context, userID string) ([]CardSummary, error) {
	return s.repo.ListUserCardSummaries(ctx, userID)
}

// BlockCard blocks a card
func (s *Service) BlockCard(ctx context.Context, userID, cardID string) error {
	// 1. Verify ownership
	card, err := s.repo.GetByIDForSummary(ctx, cardID)
	if err != nil {
		return err
	}
	if card.UserID != userID {
		return ErrUnauthorized
	}

	// 2. Check if already blocked or cancelled
	if card.Status == "blocked" {
		return nil // Already blocked
	}
	if card.Status == "cancelled" {
		return ErrCardCancelled
	}

	// 3. Update status
	return s.repo.UpdateStatus(ctx, cardID, "blocked")
}

// UnblockCard unblocks a card
func (s *Service) UnblockCard(ctx context.Context, userID, cardID string) error {
	// 1. Verify ownership
	card, err := s.repo.GetByIDForSummary(ctx, cardID)
	if err != nil {
		return err
	}
	if card.UserID != userID {
		return ErrUnauthorized
	}

	// 2. Check if card can be unblocked
	if card.Status == "cancelled" {
		return ErrCardCancelled
	}
	if card.Status != "blocked" {
		return nil // Already active
	}

	// 3. Update status
	return s.repo.UpdateStatus(ctx, cardID, "active")
}

// UpdateLimits updates spending limits
func (s *Service) UpdateLimits(ctx context.Context, userID, cardID string, req UpdateLimitsRequest) error {
	// 1. Verify ownership
	card, err := s.repo.GetByIDForSummary(ctx, cardID)
	if err != nil {
		return err
	}
	if card.UserID != userID {
		return ErrUnauthorized
	}

	// 2. Validate limits
	if req.DailyLimitCents < 0 || req.MonthlyLimitCents < 0 {
		return ErrInvalidLimit
	}

	// 3. Update limits
	return s.repo.UpdateLimits(ctx, cardID, UpdateLimitsParams{
		DailyLimitCents:   req.DailyLimitCents,
		MonthlyLimitCents: req.MonthlyLimitCents,
	})
}

// UpdateSecuritySettings updates security settings
func (s *Service) UpdateSecuritySettings(ctx context.Context, userID, cardID string, req SecuritySettingsRequest) error {
	// 1. Verify ownership
	card, err := s.repo.GetByIDForSummary(ctx, cardID)
	if err != nil {
		return err
	}
	if card.UserID != userID {
		return ErrUnauthorized
	}

	// 2. Update settings
	return s.repo.UpdateSecuritySettings(ctx, cardID, SecuritySettingsParams{
		IsContactless:      req.IsContactless,
		IsInternational:    req.IsInternational,
		BlockInternational: req.BlockInternational,
		BlockOnline:        req.BlockOnline,
	})
}

// SetPIN sets or changes card PIN
func (s *Service) SetPIN(ctx context.Context, userID, cardID string, req SetPINRequest) error {
	// 1. Validate new PIN format
	if err := ValidatePIN(req.PIN); err != nil {
		return err
	}

	// 2. Verify ownership
	card, err := s.repo.GetByIDForSummary(ctx, cardID)
	if err != nil {
		return err
	}
	if card.UserID != userID {
		return ErrUnauthorized
	}

	// 3. If changing existing PIN, verify current PIN
	if req.CurrentPIN != "" {
		match, err := s.repo.VerifyPIN(ctx, cardID, req.CurrentPIN)
		if err != nil {
			if err == ErrPINNotSet {
				// No PIN set, allow setting new PIN
			} else {
				return err
			}
		} else if !match {
			return ErrPINMismatch
		}
	}

	// 4. Update PIN
	return s.repo.UpdatePIN(ctx, cardID, req.PIN)
}

// VerifyPIN verifies a PIN (for transaction authorization)
func (s *Service) VerifyPIN(ctx context.Context, cardID, pin string) (bool, error) {
	return s.repo.VerifyPIN(ctx, cardID, pin)
}

// CancelCard cancels a card permanently
func (s *Service) CancelCard(ctx context.Context, userID, cardID string, reason string) error {
	// 1. Verify ownership
	card, err := s.repo.GetByIDForSummary(ctx, cardID)
	if err != nil {
		return err
	}
	if card.UserID != userID {
		return ErrUnauthorized
	}

	// 2. Check if already cancelled
	if card.Status == "cancelled" {
		return nil // Already cancelled
	}

	// 3. Cancel card (soft delete)
	return s.repo.CancelCard(ctx, cardID)
}

// ProcessCardTransaction processes a card transaction (checks limits, updates spent)
// This would be called by a card transaction processor
func (s *Service) ProcessCardTransaction(ctx context.Context, cardID string, amountCents int64) error {
	return s.executeInTransaction(ctx, func(tx *sql.Tx) error {
		// 1. Lock card record
		card, err := s.repo.GetForUpdate(ctx, tx, cardID)
		if err != nil {
			return err
		}

		// 2. Check card status
		if card.Status != "active" {
			if card.Status == "blocked" {
				return ErrCardBlocked
			}
			if card.Status == "cancelled" {
				return ErrCardCancelled
			}
			return ErrCardNotActive
		}

		// 3. Check if card is expired
		if card.ExpiresAt.Valid && card.ExpiresAt.Time.Before(time.Now()) {
			return ErrCardExpired
		}

		// 4. Check daily limit
		currentDailySpent := int64(0)
		if card.CurrentDailySpentCents.Valid {
			currentDailySpent = card.CurrentDailySpentCents.Int64
		}
		dailyLimit := int64(500000) // Default
		if card.DailyLimitCents.Valid {
			dailyLimit = card.DailyLimitCents.Int64
		}

		if currentDailySpent+amountCents > dailyLimit {
			return ErrDailyLimitExceeded
		}

		// 5. Check monthly limit
		currentMonthlySpent := int64(0)
		if card.CurrentMonthlySpentCents.Valid {
			currentMonthlySpent = card.CurrentMonthlySpentCents.Int64
		}
		monthlyLimit := int64(5000000) // Default
		if card.MonthlyLimitCents.Valid {
			monthlyLimit = card.MonthlyLimitCents.Int64
		}

		if currentMonthlySpent+amountCents > monthlyLimit {
			return ErrMonthlyLimitExceeded
		}

		// 6. Update spent amounts
		newDaily := currentDailySpent + amountCents
		newMonthly := currentMonthlySpent + amountCents

		return s.repo.UpdateSpentAmounts(ctx, cardID, newDaily, newMonthly)
	})
}

// executeInTransaction executes a function within a database transaction
func (s *Service) executeInTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// Ensure rollback on panic or error
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Re-throw panic after rollback
		} else if err != nil {
			tx.Rollback()
		}
	}()

	// Execute function
	err = fn(tx)
	if err != nil {
		return err
	}

	// Commit transaction
	return tx.Commit()
}
