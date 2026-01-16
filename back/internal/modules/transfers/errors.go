package transfers

import "errors"

var (
	// ErrTransferNotFound is returned when a transfer is not found
	ErrTransferNotFound = errors.New("transfer not found")

	// ErrInsufficientBalance is returned when user has insufficient balance
	ErrInsufficientBalance = errors.New("insufficient balance")

	// ErrDailyLimitExceeded is returned when daily transfer limit is exceeded
	ErrDailyLimitExceeded = errors.New("daily transfer limit exceeded")

	// ErrMonthlyLimitExceeded is returned when monthly transfer limit is exceeded
	ErrMonthlyLimitExceeded = errors.New("monthly transfer limit exceeded")

	// ErrInvalidAmount is returned when transfer amount is invalid
	ErrInvalidAmount = errors.New("invalid transfer amount")

	// ErrInvalidPIXKey is returned when PIX key format is invalid
	ErrInvalidPIXKey = errors.New("invalid PIX key")

	// ErrInvalidTransferStatus is returned when operation cannot be performed on current status
	ErrInvalidTransferStatus = errors.New("invalid transfer status for this operation")

	// ErrRecipientNotFound is returned when recipient user is not found
	ErrRecipientNotFound = errors.New("recipient user not found")

	// ErrCannotTransferToSelf is returned when trying to transfer to yourself
	ErrCannotTransferToSelf = errors.New("cannot transfer to yourself")

	// ErrInvalidBankData is returned when bank account data is invalid
	ErrInvalidBankData = errors.New("invalid bank account data")

	// ErrInvalidCPF is returned when CPF validation fails
	ErrInvalidCPF = errors.New("invalid CPF")
)
