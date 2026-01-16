package cards

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidateCardNumber validates a card number using the Luhn algorithm
func ValidateCardNumber(cardNumber string) error {
	// Remove spaces and dashes
	cleaned := regexp.MustCompile(`[\s-]`).ReplaceAllString(cardNumber, "")

	// Must be 13-19 digits
	if len(cleaned) < 13 || len(cleaned) > 19 {
		return ErrInvalidCardNumber
	}

	// Must be all digits
	if !regexp.MustCompile(`^\d+$`).MatchString(cleaned) {
		return ErrInvalidCardNumber
	}

	// Luhn algorithm validation
	if !isValidLuhn(cleaned) {
		return ErrInvalidCardNumber
	}

	return nil
}

// isValidLuhn checks if a card number passes the Luhn algorithm
func isValidLuhn(cardNumber string) bool {
	sum := 0
	alternate := false

	// Iterate from right to left
	for i := len(cardNumber) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(cardNumber[i]))

		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		alternate = !alternate
	}

	return sum%10 == 0
}

// calculateLuhnChecksum calculates the Luhn check digit for a partial card number
func calculateLuhnChecksum(partial string) int {
	sum := 0
	alternate := true // Start with alternate=true because check digit will be at the end

	// Process all existing digits from right to left
	for i := len(partial) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(partial[i]))

		if alternate {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		alternate = !alternate
	}

	// Calculate check digit
	checksum := (10 - (sum % 10)) % 10
	return checksum
}

// GenerateCardNumber generates a valid card number with Luhn checksum for the specified brand
func GenerateCardNumber(brand string) (string, error) {
	var prefix string
	var length int

	// Set IIN (Issuer Identification Number) prefix and length based on brand
	switch strings.ToLower(brand) {
	case "visa":
		prefix = "4"
		length = 16
	case "mastercard":
		// Mastercard ranges: 51-55 or 2221-2720
		// Use 5x prefix for simplicity (51-55 range)
		prefixes := []string{"51", "52", "53", "54", "55"}
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(prefixes))))
		prefix = prefixes[idx.Int64()]
		length = 16
	case "elo":
		// Elo BINs start with 636368, 438935, 504175, 451416, 636297, 506726
		// Use 636368 for simplicity
		prefix = "636368"
		length = 16
	default:
		return "", ErrInvalidCardBrand
	}

	// Generate random middle digits
	middleLength := length - len(prefix) - 1 // -1 for check digit
	if middleLength < 0 {
		return "", fmt.Errorf("invalid card number length for brand %s", brand)
	}

	middle := ""
	for i := 0; i < middleLength; i++ {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		middle += fmt.Sprint(digit)
	}

	// Combine prefix and middle digits
	partial := prefix + middle

	// Calculate Luhn check digit
	checksum := calculateLuhnChecksum(partial)

	// Generate final card number
	cardNumber := partial + fmt.Sprint(checksum)

	return cardNumber, nil
}

// ValidateCVV validates CVV format (3-4 digits)
func ValidateCVV(cvv string) error {
	if len(cvv) < 3 || len(cvv) > 4 {
		return ErrInvalidCVV
	}

	if !regexp.MustCompile(`^\d+$`).MatchString(cvv) {
		return ErrInvalidCVV
	}

	return nil
}

// ValidatePIN validates PIN format (4-6 digits) and checks for weak patterns
func ValidatePIN(pin string) error {
	// Must be 4-6 digits
	if len(pin) < 4 || len(pin) > 6 {
		return ErrInvalidPIN
	}

	// Must be all digits
	if !regexp.MustCompile(`^\d+$`).MatchString(pin) {
		return ErrInvalidPIN
	}

	// Check for weak PIN patterns
	if isWeakPIN(pin) {
		return ErrWeakPIN
	}

	return nil
}

// isWeakPIN checks for common weak PIN patterns
func isWeakPIN(pin string) bool {
	// All same digit (0000, 1111, etc.)
	if len(pin) > 0 {
		firstDigit := pin[0]
		allSame := true
		for i := 1; i < len(pin); i++ {
			if pin[i] != firstDigit {
				allSame = false
				break
			}
		}
		if allSame {
			return true
		}
	}

	// Sequential ascending (1234, 12345, etc.)
	isSequentialAsc := true
	for i := 1; i < len(pin); i++ {
		if pin[i] != pin[i-1]+1 {
			isSequentialAsc = false
			break
		}
	}
	if isSequentialAsc {
		return true
	}

	// Sequential descending (4321, 54321, etc.)
	isSequentialDesc := true
	for i := 1; i < len(pin); i++ {
		if pin[i] != pin[i-1]-1 {
			isSequentialDesc = false
			break
		}
	}
	if isSequentialDesc {
		return true
	}

	// Common weak PINs
	weakPINs := map[string]bool{
		"1234":   true,
		"4321":   true,
		"1111":   true,
		"2222":   true,
		"0000":   true,
		"123456": true,
		"654321": true,
	}

	if weakPINs[pin] {
		return true
	}

	return false
}

// ValidateExpiryDate validates expiry month and year
func ValidateExpiryDate(month, year int) error {
	// Validate month range
	if month < 1 || month > 12 {
		return ErrInvalidExpiryDate
	}

	// Validate year is not in the past
	currentYear := time.Now().Year()
	if year < currentYear {
		return ErrInvalidExpiryDate
	}

	// If it's the current year, validate month is not in the past
	if year == currentYear {
		currentMonth := int(time.Now().Month())
		if month < currentMonth {
			return ErrInvalidExpiryDate
		}
	}

	return nil
}

// CalculateExpiryDate calculates the expiry date for a card based on type
// Physical cards: 5 years from now
// Virtual cards: 3 years from now
func CalculateExpiryDate(cardType string) time.Time {
	now := time.Now()

	var years int
	switch strings.ToLower(cardType) {
	case "physical":
		years = 5
	case "virtual":
		years = 3
	default:
		years = 3 // Default to 3 years
	}

	return now.AddDate(years, 0, 0)
}

// ValidateCardType validates card type
func ValidateCardType(cardType string) error {
	validTypes := map[string]bool{
		"physical": true,
		"virtual":  true,
	}

	if !validTypes[strings.ToLower(cardType)] {
		return ErrInvalidCardType
	}

	return nil
}

// ValidateCardBrand validates card brand
func ValidateCardBrand(brand string) error {
	validBrands := map[string]bool{
		"visa":       true,
		"mastercard": true,
		"elo":        true,
	}

	if !validBrands[strings.ToLower(brand)] {
		return ErrInvalidCardBrand
	}

	return nil
}
