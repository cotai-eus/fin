package bills

import "strings"

// ValidateBillType validates bill type
func ValidateBillType(billType string) error {
	validTypes := map[string]bool{
		"utility":     true,
		"internet":    true,
		"phone":       true,
		"credit_card": true,
		"insurance":   true,
		"other":       true,
	}

	if !validTypes[billType] {
		return ErrInvalidType
	}
	return nil
}

// ValidateAmount validates bill amount
func ValidateAmount(amount int64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	// Max amount: R$ 100,000.00
	if amount > 10000000 {
		return ErrInvalidAmount
	}
	return nil
}

// NormalizeBarcode removes spaces and dashes from barcode
func NormalizeBarcode(barcode string) string {
	barcode = strings.ReplaceAll(barcode, " ", "")
	barcode = strings.ReplaceAll(barcode, "-", "")
	barcode = strings.ReplaceAll(barcode, ".", "")
	return barcode
}
