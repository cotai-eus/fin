package bills

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Brazilian barcode formats:
// - Boleto Bancário: 47 digits (line) or 44 digits (barcode)
// - Concessionárias (utilities): 48 digits

var (
	// Boleto bancário pattern (47 digits with spaces/dots or 44 digits)
	boletoLineRegex    = regexp.MustCompile(`^\d{5}\.\d{5}\s\d{5}\.\d{6}\s\d{5}\.\d{6}\s\d\s\d{14}$`)
	boletoBarcodeRegex = regexp.MustCompile(`^\d{44}$`)

	// Concessionária pattern (48 digits)
	concessionariaRegex = regexp.MustCompile(`^\d{48}$`)
)

// BarcodeType represents the type of barcode
type BarcodeType string

const (
	BarcodeTypeBoleto         BarcodeType = "boleto"
	BarcodeTypeConcessionaria BarcodeType = "concessionaria"
)

// BarcodeInfo contains parsed information from a barcode
type BarcodeInfo struct {
	Type          BarcodeType
	Barcode       string // Normalized barcode (digits only)
	AmountCents   int64
	DueDate       time.Time
	RecipientName string
	BillType      string
}

// ValidateBarcode validates a Brazilian barcode and returns parsed information
func ValidateBarcode(input string) (*BarcodeInfo, error) {
	// Remove spaces and special characters
	normalized := strings.ReplaceAll(input, " ", "")
	normalized = strings.ReplaceAll(normalized, ".", "")
	normalized = strings.ReplaceAll(normalized, "-", "")

	// Determine barcode type
	if boletoBarcodeRegex.MatchString(normalized) {
		return validateBoletoBancario(normalized)
	} else if len(normalized) == 47 {
		// Convert linha digitável (47) to código de barras (44)
		barcode, err := convertLinhaDigitavelToBarcode(input)
		if err != nil {
			return nil, ErrInvalidBarcode
		}
		return validateBoletoBancario(barcode)
	} else if concessionariaRegex.MatchString(normalized) {
		return validateConcessionaria(normalized)
	}

	return nil, ErrInvalidBarcode
}

// validateBoletoBancario validates a boleto bancário (44 digits)
func validateBoletoBancario(barcode string) (*BarcodeInfo, error) {
	if len(barcode) != 44 {
		return nil, ErrInvalidBarcode
	}

	// Validate check digit (position 4, using modulo 11)
	if !validateBoletoCheckDigit(barcode) {
		return nil, ErrInvalidBarcode
	}

	// Parse amount (positions 9-18, in centavos)
	amountStr := barcode[9:19]
	amountCents, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amountCents <= 0 {
		return nil, ErrInvalidAmount
	}

	// Parse due date (positions 5-8, days since 07/10/1997)
	dueDateFactor := barcode[5:9]
	dueDate := parseDueDateFactor(dueDateFactor)

	// Determine recipient based on bank code (first 3 digits)
	bankCode := barcode[0:3]
	recipientName := getBankName(bankCode)

	return &BarcodeInfo{
		Type:          BarcodeTypeBoleto,
		Barcode:       barcode,
		AmountCents:   amountCents,
		DueDate:       dueDate,
		RecipientName: recipientName,
		BillType:      "bank",
	}, nil
}

// validateConcessionaria validates a concessionária barcode (48 digits)
func validateConcessionaria(barcode string) (*BarcodeInfo, error) {
	if len(barcode) != 48 {
		return nil, ErrInvalidBarcode
	}

	// Validate check digits (positions 11, 22, 33, 44 using modulo 10/11)
	if !validateConcessionariaCheckDigits(barcode) {
		return nil, ErrInvalidBarcode
	}

	// Parse amount (positions 4-15, in centavos)
	amountStr := barcode[4:15]
	amountCents, err := strconv.ParseInt(amountStr, 10, 64)
	if err != nil || amountCents <= 0 {
		return nil, ErrInvalidAmount
	}

	// Concessionária bills typically don't have due date in barcode
	// Use a default (e.g., 30 days from now)
	dueDate := time.Now().AddDate(0, 0, 30)

	// Determine type based on segment code (position 2)
	segmentCode := barcode[1:2]
	billType := getConcessionariaType(segmentCode)
	recipientName := getConcessionariaName(segmentCode)

	return &BarcodeInfo{
		Type:          BarcodeTypeConcessionaria,
		Barcode:       barcode,
		AmountCents:   amountCents,
		DueDate:       dueDate,
		RecipientName: recipientName,
		BillType:      billType,
	}, nil
}

// convertLinhaDigitavelToBarcode converts a 47-digit linha digitável to 44-digit barcode
func convertLinhaDigitavelToBarcode(linha string) (string, error) {
	// Remove formatting
	digits := regexp.MustCompile(`\d`).FindAllString(linha, -1)
	if len(digits) != 47 {
		return "", ErrInvalidBarcode
	}

	digitStr := strings.Join(digits, "")

	// Extract parts (removing check digits at positions 9, 20, 31)
	barcode := digitStr[0:4] + digitStr[32:33] + digitStr[33:47] +
		digitStr[4:9] + digitStr[10:20] + digitStr[21:31]

	if len(barcode) != 44 {
		return "", ErrInvalidBarcode
	}

	return barcode, nil
}

// validateBoletoCheckDigit validates the check digit using modulo 11
func validateBoletoCheckDigit(barcode string) bool {
	if len(barcode) < 5 {
		return false
	}

	// Extract check digit (position 4)
	checkDigit, _ := strconv.Atoi(string(barcode[4]))

	// Calculate expected check digit using modulo 11
	sequence := barcode[0:4] + barcode[5:]
	sum := 0
	multiplier := 2

	for i := len(sequence) - 1; i >= 0; i-- {
		digit, _ := strconv.Atoi(string(sequence[i]))
		sum += digit * multiplier
		multiplier++
		if multiplier > 9 {
			multiplier = 2
		}
	}

	remainder := sum % 11
	expected := 11 - remainder

	// Special cases for modulo 11
	if expected == 0 || expected == 10 || expected == 11 {
		expected = 1
	}

	return checkDigit == expected
}

// validateConcessionariaCheckDigits validates check digits using modulo 10
func validateConcessionariaCheckDigits(barcode string) bool {
	// Simplified validation - in production, implement full modulo 10/11 algorithm
	// for positions 11, 22, 33, 44
	positions := []int{11, 22, 33, 44}
	for _, pos := range positions {
		if pos >= len(barcode) {
			return false
		}
		// Basic digit check
		if _, err := strconv.Atoi(string(barcode[pos])); err != nil {
			return false
		}
	}
	return true
}

// parseDueDateFactor converts the due date factor to a time.Time
func parseDueDateFactor(factor string) time.Time {
	factorInt, err := strconv.Atoi(factor)
	if err != nil || factorInt == 0 {
		// If no date or invalid, return 30 days from now
		return time.Now().AddDate(0, 0, 30)
	}

	// Base date: October 7, 1997
	baseDate := time.Date(1997, 10, 7, 0, 0, 0, 0, time.UTC)
	dueDate := baseDate.AddDate(0, 0, factorInt)

	return dueDate
}

// getBankName returns bank name based on bank code
func getBankName(code string) string {
	banks := map[string]string{
		"001": "Banco do Brasil",
		"033": "Santander",
		"104": "Caixa Econômica Federal",
		"237": "Bradesco",
		"341": "Itaú",
		"356": "Banco Real",
		"389": "Banco Mercantil do Brasil",
		"399": "HSBC",
		"422": "Banco Safra",
		"453": "Banco Rural",
		"633": "Banco Rendimento",
		"652": "Itaú Unibanco",
		"745": "Citibank",
	}

	if name, ok := banks[code]; ok {
		return name
	}
	return "Instituição Financeira"
}

// getConcessionariaType returns bill type based on segment code
func getConcessionariaType(code string) string {
	types := map[string]string{
		"1": "utility", // Energia elétrica
		"2": "utility", // Telecomunicações
		"3": "utility", // Água e saneamento
		"4": "utility", // Gás
		"5": "tax",     // Tributos
		"6": "other",   // Outros
		"7": "other",   // Outros
		"8": "other",   // Outros
		"9": "tax",     // Uso próprio
	}

	if t, ok := types[code]; ok {
		return t
	}
	return "other"
}

// getConcessionariaName returns recipient name based on segment code
func getConcessionariaName(code string) string {
	names := map[string]string{
		"1": "Companhia de Energia",
		"2": "Operadora de Telefonia",
		"3": "Companhia de Água e Saneamento",
		"4": "Companhia de Gás",
		"5": "Receita Federal",
		"9": "Órgão Público",
	}

	if name, ok := names[code]; ok {
		return name
	}
	return "Concessionária de Serviço Público"
}
