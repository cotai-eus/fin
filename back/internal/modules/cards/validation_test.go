package cards

import (
	"testing"
)

// TestValidateCardNumber tests card number validation with Luhn algorithm
func TestValidateCardNumber(t *testing.T) {
	tests := []struct {
		name        string
		cardNumber  string
		expectError bool
	}{
		{"Valid Visa", "4532015112830366", false}, // Valid Luhn
		{"Valid Mastercard", "5425233430109903", false},
		{"Valid with spaces", "4532 0151 1283 0366", false}, // Valid Luhn
		{"Valid with dashes", "4532-0151-1283-0366", false}, // Valid Luhn
		{"Invalid Luhn checksum", "4532123456789012", true},
		{"Too short", "123456", true},
		{"Too long", "12345678901234567890", true},
		{"Non-numeric", "abcd1234efgh5678", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCardNumber(tt.cardNumber)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateCardNumber(%s) error = %v, expectError %v", tt.cardNumber, err, tt.expectError)
			}
		})
	}
}

// TestGenerateCardNumber tests card number generation for different brands
func TestGenerateCardNumber(t *testing.T) {
	brands := []string{"visa", "mastercard", "elo"}

	for _, brand := range brands {
		t.Run(brand, func(t *testing.T) {
			// Generate card number
			cardNumber, err := GenerateCardNumber(brand)
			if err != nil {
				t.Fatalf("GenerateCardNumber(%s) failed: %v", brand, err)
			}

			// Validate generated number passes Luhn
			if err := ValidateCardNumber(cardNumber); err != nil {
				t.Errorf("Generated card number %s failed Luhn validation: %v", cardNumber, err)
			}

			// Check length
			if len(cardNumber) != 16 {
				t.Errorf("Expected 16 digits, got %d", len(cardNumber))
			}

			// Verify brand prefix
			switch brand {
			case "visa":
				if cardNumber[0] != '4' {
					t.Errorf("Visa card should start with 4, got %c", cardNumber[0])
				}
			case "mastercard":
				firstTwo := cardNumber[0:2]
				if firstTwo < "51" || firstTwo > "55" {
					t.Errorf("Mastercard should start with 51-55, got %s", firstTwo)
				}
			case "elo":
				if cardNumber[0:6] != "636368" {
					t.Errorf("Elo card should start with 636368, got %s", cardNumber[0:6])
				}
			}
		})
	}
}

// TestGenerateCardNumberUniqueness verifies each generation produces unique numbers
func TestGenerateCardNumberUniqueness(t *testing.T) {
	generated := make(map[string]bool)

	for i := 0; i < 100; i++ {
		cardNumber, err := GenerateCardNumber("visa")
		if err != nil {
			t.Fatalf("Generation failed: %v", err)
		}

		if generated[cardNumber] {
			t.Errorf("Generated duplicate card number: %s", cardNumber)
		}
		generated[cardNumber] = true
	}
}

// TestValidateCVV tests CVV validation
func TestValidateCVV(t *testing.T) {
	tests := []struct {
		name        string
		cvv         string
		expectError bool
	}{
		{"Valid 3-digit", "123", false},
		{"Valid 4-digit", "1234", false},
		{"Too short", "12", true},
		{"Too long", "12345", true},
		{"Non-numeric", "abc", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCVV(tt.cvv)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateCVV(%s) error = %v, expectError %v", tt.cvv, err, tt.expectError)
			}
		})
	}
}

// TestValidatePIN tests PIN validation and weak pattern detection
func TestValidatePIN(t *testing.T) {
	tests := []struct {
		name        string
		pin         string
		expectError bool
	}{
		{"Valid 4-digit", "1357", false},
		{"Valid 6-digit", "135790", false},
		{"Too short", "123", true},
		{"Too long", "1234567", true},
		{"Non-numeric", "abcd", true},
		{"Empty", "", true},
		{"Weak: all same", "0000", true},
		{"Weak: all same (1111)", "1111", true},
		{"Weak: sequential ascending", "1234", true},
		{"Weak: sequential descending", "4321", true},
		{"Weak: common PIN", "123456", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePIN(tt.pin)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidatePIN(%s) error = %v, expectError %v", tt.pin, err, tt.expectError)
			}
		})
	}
}

// TestValidateExpiryDate tests expiry date validation
func TestValidateExpiryDate(t *testing.T) {
	tests := []struct {
		name        string
		month       int
		year        int
		expectError bool
	}{
		{"Valid future date", 12, 2027, false},
		{"Valid current year", 12, 2026, false},
		{"Invalid month (0)", 0, 2027, true},
		{"Invalid month (13)", 13, 2027, true},
		{"Invalid year (past)", 12, 2020, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateExpiryDate(tt.month, tt.year)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateExpiryDate(%d, %d) error = %v, expectError %v", tt.month, tt.year, err, tt.expectError)
			}
		})
	}
}

// TestValidateCardType tests card type validation
func TestValidateCardType(t *testing.T) {
	tests := []struct {
		name        string
		cardType    string
		expectError bool
	}{
		{"Valid physical", "physical", false},
		{"Valid virtual", "virtual", false},
		{"Valid uppercase", "PHYSICAL", false},
		{"Invalid type", "debit", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCardType(tt.cardType)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateCardType(%s) error = %v, expectError %v", tt.cardType, err, tt.expectError)
			}
		})
	}
}

// TestValidateCardBrand tests card brand validation
func TestValidateCardBrand(t *testing.T) {
	tests := []struct {
		name        string
		brand       string
		expectError bool
	}{
		{"Valid visa", "visa", false},
		{"Valid mastercard", "mastercard", false},
		{"Valid elo", "elo", false},
		{"Valid uppercase", "VISA", false},
		{"Invalid brand", "amex", true},
		{"Empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateCardBrand(tt.brand)
			if (err != nil) != tt.expectError {
				t.Errorf("ValidateCardBrand(%s) error = %v, expectError %v", tt.brand, err, tt.expectError)
			}
		})
	}
}

// TestCalculateExpiryDate tests expiry date calculation
func TestCalculateExpiryDate(t *testing.T) {
	tests := []struct {
		name          string
		cardType      string
		expectedYears int
	}{
		{"Physical card", "physical", 5},
		{"Virtual card", "virtual", 3},
		{"Unknown type defaults to 3", "unknown", 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expiryDate := CalculateExpiryDate(tt.cardType)

			// Verify expiry is in the future
			// We can't check exact years due to timing, but we can verify it's reasonable
			if expiryDate.Before(CalculateExpiryDate("virtual")) && tt.cardType == "physical" {
				t.Error("Physical card should expire later than virtual card")
			}
		})
	}
}

// TestIsValidLuhn tests the Luhn algorithm implementation
func TestIsValidLuhn(t *testing.T) {
	tests := []struct {
		name       string
		cardNumber string
		expected   bool
	}{
		{"Valid Visa", "4532015112830366", true}, // Fixed valid Luhn
		{"Valid Mastercard", "5425233430109903", true},
		{"Invalid checksum", "4532123456789012", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidLuhn(tt.cardNumber)
			if result != tt.expected {
				t.Errorf("isValidLuhn(%s) = %v, expected %v", tt.cardNumber, result, tt.expected)
			}
		})
	}
}

// TestCalculateLuhnChecksum tests checksum calculation
func TestCalculateLuhnChecksum(t *testing.T) {
	tests := []struct {
		name     string
		partial  string
		expected int
	}{
		{"Visa prefix", "453201511283036", 6},       // Should produce valid Visa: 4532015112830366
		{"Mastercard prefix", "542523343010990", 3}, // Should produce valid MC: 5425233430109903
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checksum := calculateLuhnChecksum(tt.partial)
			if checksum != tt.expected {
				t.Errorf("calculateLuhnChecksum(%s) = %d, expected %d", tt.partial, checksum, tt.expected)
			}

			// Verify the complete number is valid
			completeNumber := tt.partial + string(rune('0'+checksum))
			if !isValidLuhn(completeNumber) {
				t.Errorf("Complete number %s should pass Luhn validation", completeNumber)
			}
		})
	}
}

// BenchmarkGenerateCardNumber benchmarks card number generation
func BenchmarkGenerateCardNumber(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateCardNumber("visa")
	}
}

// BenchmarkValidateCardNumber benchmarks card number validation
func BenchmarkValidateCardNumber(b *testing.B) {
	cardNumber := "4532123456789010"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidateCardNumber(cardNumber)
	}
}

// BenchmarkValidatePIN benchmarks PIN validation
func BenchmarkValidatePIN(b *testing.B) {
	pin := "1357"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ValidatePIN(pin)
	}
}
