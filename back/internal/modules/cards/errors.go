package cards

import "errors"

var (
	// Card errors
	ErrCardNotFound  = errors.New("card not found")
	ErrCardNotActive = errors.New("card is not active")
	ErrCardExpired   = errors.New("card has expired")
	ErrCardCancelled = errors.New("card has been cancelled")
	ErrCardBlocked   = errors.New("card is blocked")

	// Validation errors
	ErrInvalidCardNumber = errors.New("invalid card number")
	ErrInvalidCVV        = errors.New("invalid CVV")
	ErrInvalidPIN        = errors.New("invalid PIN format")
	ErrWeakPIN           = errors.New("weak PIN detected")
	ErrInvalidExpiryDate = errors.New("invalid expiry date")
	ErrInvalidCardType   = errors.New("invalid card type")
	ErrInvalidCardBrand  = errors.New("invalid card brand")

	// Limit errors
	ErrDailyLimitExceeded   = errors.New("daily spending limit exceeded")
	ErrMonthlyLimitExceeded = errors.New("monthly spending limit exceeded")
	ErrInvalidLimit         = errors.New("invalid limit amount")

	// PIN errors
	ErrPINNotSet    = errors.New("PIN not set for this card")
	ErrPINIncorrect = errors.New("incorrect PIN")
	ErrPINMismatch  = errors.New("current PIN does not match")

	// Security errors
	ErrInternationalBlocked = errors.New("international transactions blocked")
	ErrOnlineBlocked        = errors.New("online transactions blocked")
	ErrContactlessBlocked   = errors.New("contactless transactions blocked")

	// Encryption errors
	ErrEncryptionFailed = errors.New("encryption failed")
	ErrDecryptionFailed = errors.New("decryption failed")

	// Authorization errors
	ErrUnauthorized = errors.New("unauthorized access to card")
)
