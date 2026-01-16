package transfers

import (
	"net/mail"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// ValidateAmount validates that the transfer amount is positive
func ValidateAmount(amountCents int64) error {
	if amountCents <= 0 {
		return ErrInvalidAmount
	}
	return nil
}

// ValidatePIXKey validates a PIX key based on its type
func ValidatePIXKey(key, keyType string) error {
	if key == "" || keyType == "" {
		return ErrInvalidPIXKey
	}

	switch keyType {
	case "cpf":
		return ValidateCPF(key)
	case "cnpj":
		return ValidateCNPJ(key)
	case "email":
		return ValidateEmail(key)
	case "phone":
		return ValidatePhone(key)
	case "random":
		return ValidateRandomKey(key)
	default:
		return ErrInvalidPIXKey
	}
}

// ValidateCPF validates a Brazilian CPF (11 digits)
func ValidateCPF(cpf string) error {
	// Remove non-numeric characters
	cpf = regexp.MustCompile(`\D`).ReplaceAllString(cpf, "")

	// Must have exactly 11 digits
	if len(cpf) != 11 {
		return ErrInvalidCPF
	}

	// Check for known invalid CPFs (all same digit)
	if cpf == "00000000000" || cpf == "11111111111" || cpf == "22222222222" ||
		cpf == "33333333333" || cpf == "44444444444" || cpf == "55555555555" ||
		cpf == "66666666666" || cpf == "77777777777" || cpf == "88888888888" ||
		cpf == "99999999999" {
		return ErrInvalidCPF
	}

	// Validate check digits
	var sum int
	var weight int

	// First check digit
	weight = 10
	for i := 0; i < 9; i++ {
		num, _ := strconv.Atoi(string(cpf[i]))
		sum += num * weight
		weight--
	}

	remainder := sum % 11
	checkDigit1 := 0
	if remainder >= 2 {
		checkDigit1 = 11 - remainder
	}

	firstDigit, _ := strconv.Atoi(string(cpf[9]))
	if firstDigit != checkDigit1 {
		return ErrInvalidCPF
	}

	// Second check digit
	sum = 0
	weight = 11
	for i := 0; i < 10; i++ {
		num, _ := strconv.Atoi(string(cpf[i]))
		sum += num * weight
		weight--
	}

	remainder = sum % 11
	checkDigit2 := 0
	if remainder >= 2 {
		checkDigit2 = 11 - remainder
	}

	secondDigit, _ := strconv.Atoi(string(cpf[10]))
	if secondDigit != checkDigit2 {
		return ErrInvalidCPF
	}

	return nil
}

// ValidateCNPJ validates a Brazilian CNPJ (14 digits)
func ValidateCNPJ(cnpj string) error {
	// Remove non-numeric characters
	cnpj = regexp.MustCompile(`\D`).ReplaceAllString(cnpj, "")

	// Must have exactly 14 digits
	if len(cnpj) != 14 {
		return ErrInvalidPIXKey
	}

	// Check for known invalid CNPJs (all same digit)
	if cnpj == "00000000000000" || cnpj == "11111111111111" || cnpj == "22222222222222" ||
		cnpj == "33333333333333" || cnpj == "44444444444444" || cnpj == "55555555555555" ||
		cnpj == "66666666666666" || cnpj == "77777777777777" || cnpj == "88888888888888" ||
		cnpj == "99999999999999" {
		return ErrInvalidPIXKey
	}

	// Validate check digits
	weights1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	weights2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}

	// First check digit
	sum := 0
	for i := 0; i < 12; i++ {
		num, _ := strconv.Atoi(string(cnpj[i]))
		sum += num * weights1[i]
	}
	remainder := sum % 11
	checkDigit1 := 0
	if remainder >= 2 {
		checkDigit1 = 11 - remainder
	}

	firstDigit, _ := strconv.Atoi(string(cnpj[12]))
	if firstDigit != checkDigit1 {
		return ErrInvalidPIXKey
	}

	// Second check digit
	sum = 0
	for i := 0; i < 13; i++ {
		num, _ := strconv.Atoi(string(cnpj[i]))
		sum += num * weights2[i]
	}
	remainder = sum % 11
	checkDigit2 := 0
	if remainder >= 2 {
		checkDigit2 = 11 - remainder
	}

	secondDigit, _ := strconv.Atoi(string(cnpj[13]))
	if secondDigit != checkDigit2 {
		return ErrInvalidPIXKey
	}

	return nil
}

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrInvalidPIXKey
	}
	return nil
}

// ValidatePhone validates a Brazilian phone number
// Expected format: +55XXXXXXXXXXX or 55XXXXXXXXXXX (country code + 10-11 digits)
func ValidatePhone(phone string) error {
	// Remove non-numeric characters except +
	phone = strings.TrimSpace(phone)
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Remove leading + if present
	phone = strings.TrimPrefix(phone, "+")

	// Must start with 55 (Brazil country code) and have 12-13 digits total
	if !strings.HasPrefix(phone, "55") {
		return ErrInvalidPIXKey
	}

	// Remove country code for length validation
	phoneWithoutCountry := phone[2:]

	// Must have 10 or 11 digits (area code + number)
	if len(phoneWithoutCountry) < 10 || len(phoneWithoutCountry) > 11 {
		return ErrInvalidPIXKey
	}

	// Check if all characters are digits
	matched, _ := regexp.MatchString(`^\d+$`, phone)
	if !matched {
		return ErrInvalidPIXKey
	}

	return nil
}

// ValidateRandomKey validates a random PIX key (UUID format)
func ValidateRandomKey(key string) error {
	_, err := uuid.Parse(key)
	if err != nil {
		return ErrInvalidPIXKey
	}
	return nil
}

// ValidateTEDData validates TED transfer data
func ValidateTEDData(req CreateTEDRequest) error {
	// Validate recipient name
	if strings.TrimSpace(req.RecipientName) == "" {
		return ErrInvalidBankData
	}

	// Validate CPF/CNPJ
	if err := ValidateCPF(req.RecipientDocument); err != nil {
		// Try CNPJ if CPF fails
		if err := ValidateCNPJ(req.RecipientDocument); err != nil {
			return ErrInvalidBankData
		}
	}

	// Validate bank code (3 digits)
	matched, _ := regexp.MatchString(`^\d{3}$`, req.RecipientBank)
	if !matched {
		return ErrInvalidBankData
	}

	// Validate branch (4-5 digits, may include check digit)
	matched, _ = regexp.MatchString(`^\d{4,5}$`, regexp.MustCompile(`\D`).ReplaceAllString(req.RecipientBranch, ""))
	if !matched {
		return ErrInvalidBankData
	}

	// Validate account (up to 12 digits, may include check digit)
	accountDigits := regexp.MustCompile(`\D`).ReplaceAllString(req.RecipientAccount, "")
	if len(accountDigits) == 0 || len(accountDigits) > 12 {
		return ErrInvalidBankData
	}

	// Validate account type
	if req.RecipientAccountType != "checking" && req.RecipientAccountType != "savings" {
		return ErrInvalidBankData
	}

	// Validate amount
	if err := ValidateAmount(req.AmountCents); err != nil {
		return err
	}

	return nil
}
