package bills

import "errors"

var (
	// ErrBillNotFound is returned when a bill is not found
	ErrBillNotFound = errors.New("bill not found")

	// ErrInvalidBarcode is returned when barcode format is invalid
	ErrInvalidBarcode = errors.New("invalid barcode format")

	// ErrBarcodeAlreadyExists is returned when barcode already registered
	ErrBarcodeAlreadyExists = errors.New("barcode already registered")

	// ErrBillAlreadyPaid is returned when trying to pay an already paid bill
	ErrBillAlreadyPaid = errors.New("bill already paid")

	// ErrBillExpired is returned when bill is past due date
	ErrBillExpired = errors.New("bill is expired")

	// ErrBillCancelled is returned when operating on a cancelled bill
	ErrBillCancelled = errors.New("bill is cancelled")

	// ErrInsufficientBalance is returned when user has insufficient balance
	ErrInsufficientBalance = errors.New("insufficient balance")

	// ErrInvalidAmount is returned when bill amount is invalid
	ErrInvalidAmount = errors.New("invalid bill amount")

	// ErrInvalidType is returned when bill type is invalid
	ErrInvalidType = errors.New("invalid bill type")

	// ErrUnauthorized is returned when user doesn't own the bill
	ErrUnauthorized = errors.New("unauthorized to access this bill")
)
