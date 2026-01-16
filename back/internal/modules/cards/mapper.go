package cards

import (
	"fmt"
	"regexp"
	"strings"

	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

// dbCardToCard converts database card to domain card (with decrypted data)
// This function should only be used internally when decrypted card data is needed
func dbCardToCard(dbCard *db.Card, cardNumber, cvv string) *Card {
	card := &Card{
		ID:                       dbCard.ID.String(),
		UserID:                   dbCard.UserID.String(),
		Type:                     dbCard.Type,
		Brand:                    dbCard.Brand,
		Status:                   dbCard.Status,
		CardNumber:               cardNumber,
		CVV:                      cvv,
		LastFourDigits:           dbCard.LastFourDigits,
		HolderName:               dbCard.HolderName,
		ExpiryMonth:              int(dbCard.ExpiryMonth),
		ExpiryYear:               int(dbCard.ExpiryYear),
		CurrentDailySpentCents:   dbCard.CurrentDailySpentCents.Int64,
		CurrentMonthlySpentCents: dbCard.CurrentMonthlySpentCents.Int64,
		CreatedAt:                dbCard.CreatedAt.Time,
		UpdatedAt:                dbCard.UpdatedAt.Time,
	}

	// Handle nullable fields
	if dbCard.DailyLimitCents.Valid {
		card.DailyLimitCents = dbCard.DailyLimitCents.Int64
	}
	if dbCard.MonthlyLimitCents.Valid {
		card.MonthlyLimitCents = dbCard.MonthlyLimitCents.Int64
	}
	if dbCard.IsContactless.Valid {
		card.IsContactless = dbCard.IsContactless.Bool
	}
	if dbCard.IsInternational.Valid {
		card.IsInternational = dbCard.IsInternational.Bool
	}
	if dbCard.BlockInternational.Valid {
		card.BlockInternational = dbCard.BlockInternational.Bool
	}
	if dbCard.BlockOnline.Valid {
		card.BlockOnline = dbCard.BlockOnline.Bool
	}
	if dbCard.ExpiresAt.Valid {
		card.ExpiresAt = dbCard.ExpiresAt.Time
	}
	if dbCard.BlockedAt.Valid {
		card.BlockedAt = &dbCard.BlockedAt.Time
	}

	return card
}

// dbCardToCardSummary converts database card to summary (no sensitive data)
// Use this for list responses where sensitive data should not be included
func dbCardToCardSummary(dbCard *db.Card) *CardSummary {
	summary := &CardSummary{
		ID:                       dbCard.ID.String(),
		UserID:                   dbCard.UserID.String(),
		Type:                     dbCard.Type,
		Brand:                    dbCard.Brand,
		Status:                   dbCard.Status,
		LastFourDigits:           dbCard.LastFourDigits,
		HolderName:               dbCard.HolderName,
		ExpiryMonth:              int(dbCard.ExpiryMonth),
		ExpiryYear:               int(dbCard.ExpiryYear),
		CurrentDailySpentCents:   dbCard.CurrentDailySpentCents.Int64,
		CurrentMonthlySpentCents: dbCard.CurrentMonthlySpentCents.Int64,
		CreatedAt:                dbCard.CreatedAt.Time,
	}

	// Handle nullable fields
	if dbCard.DailyLimitCents.Valid {
		summary.DailyLimitCents = dbCard.DailyLimitCents.Int64
	}
	if dbCard.MonthlyLimitCents.Valid {
		summary.MonthlyLimitCents = dbCard.MonthlyLimitCents.Int64
	}
	if dbCard.IsContactless.Valid {
		summary.IsContactless = dbCard.IsContactless.Bool
	}
	if dbCard.IsInternational.Valid {
		summary.IsInternational = dbCard.IsInternational.Bool
	}

	return summary
}

// cardToCardDetails converts card to details (with masked number)
// Use this for individual card detail responses (GET /cards/{id})
func cardToCardDetails(card *Card, hasPIN bool) *CardDetails {
	details := &CardDetails{
		ID:                       card.ID,
		Type:                     card.Type,
		Brand:                    card.Brand,
		Status:                   card.Status,
		LastFourDigits:           card.LastFourDigits,
		MaskedCardNumber:         maskCardNumber(card.CardNumber),
		HolderName:               card.HolderName,
		ExpiryMonth:              card.ExpiryMonth,
		ExpiryYear:               card.ExpiryYear,
		DailyLimitCents:          card.DailyLimitCents,
		MonthlyLimitCents:        card.MonthlyLimitCents,
		CurrentDailySpentCents:   card.CurrentDailySpentCents,
		CurrentMonthlySpentCents: card.CurrentMonthlySpentCents,
		IsContactless:            card.IsContactless,
		IsInternational:          card.IsInternational,
		BlockInternational:       card.BlockInternational,
		BlockOnline:              card.BlockOnline,
		HasPIN:                   hasPIN,
		CreatedAt:                card.CreatedAt,
		UpdatedAt:                card.UpdatedAt,
		ExpiresAt:                card.ExpiresAt,
		BlockedAt:                card.BlockedAt,
	}

	return details
}

// maskCardNumber masks a card number for display
// Example: "4532123456789012" -> "**** **** **** 9012"
func maskCardNumber(cardNumber string) string {
	if len(cardNumber) < 4 {
		return "****"
	}

	// Get last 4 digits
	lastFour := cardNumber[len(cardNumber)-4:]

	// Determine how many groups of 4 digits we need
	// Most cards are 16 digits (4 groups), but some are 13-19 digits
	numDigits := len(cardNumber)
	numGroups := (numDigits + 3) / 4 // Round up division

	// Create masked groups
	masked := make([]string, numGroups)
	for i := 0; i < numGroups-1; i++ {
		masked[i] = "****"
	}
	masked[numGroups-1] = lastFour

	return strings.Join(masked, " ")
}

// dbCardsToCardSummaries converts multiple database cards to summaries
func dbCardsToCardSummaries(dbCards []db.Card) []CardSummary {
	summaries := make([]CardSummary, len(dbCards))
	for i, dbCard := range dbCards {
		summaries[i] = *dbCardToCardSummary(&dbCard)
	}
	return summaries
}

// formatCardNumber formats a card number with spaces for display
// Example: "4532123456789012" -> "4532 1234 5678 9012"
func formatCardNumber(cardNumber string) string {
	// Remove any existing spaces
	cleaned := strings.ReplaceAll(cardNumber, " ", "")

	// Insert spaces every 4 digits
	var formatted string
	for i, char := range cleaned {
		if i > 0 && i%4 == 0 {
			formatted += " "
		}
		formatted += string(char)
	}

	return formatted
}

// cleanCardNumber removes spaces and dashes from a card number
func cleanCardNumber(cardNumber string) string {
	cleaned := regexp.MustCompile(`[\s-]`).ReplaceAllString(cardNumber, "")
	return cleaned
}

// getCardBrandFromNumber determines the card brand from the card number
// This is useful for validation when brand is not provided
func getCardBrandFromNumber(cardNumber string) string {
	if len(cardNumber) == 0 {
		return ""
	}

	// Remove spaces and dashes
	cleaned := cleanCardNumber(cardNumber)

	// Visa: starts with 4
	if cleaned[0] == '4' {
		return "visa"
	}

	// Mastercard: starts with 51-55 or 2221-2720
	if len(cleaned) >= 2 {
		firstTwo := cleaned[0:2]
		if firstTwo >= "51" && firstTwo <= "55" {
			return "mastercard"
		}
	}
	if len(cleaned) >= 4 {
		firstFour := cleaned[0:4]
		if firstFour >= "2221" && firstFour <= "2720" {
			return "mastercard"
		}
	}

	// Elo: starts with 636368, 438935, 504175, 451416, 636297, 506726, 650
	if len(cleaned) >= 6 {
		firstSix := cleaned[0:6]
		eloPrefixes := []string{"636368", "438935", "504175", "451416", "636297", "506726"}
		for _, prefix := range eloPrefixes {
			if firstSix == prefix {
				return "elo"
			}
		}
	}
	if len(cleaned) >= 3 {
		firstThree := cleaned[0:3]
		if firstThree == "650" {
			return "elo"
		}
	}

	return "unknown"
}

// formatExpiryDate formats expiry month and year for display
// Example: (3, 2027) -> "03/27"
func formatExpiryDate(month, year int) string {
	// Get last 2 digits of year
	yearShort := year % 100
	return fmt.Sprintf("%02d/%02d", month, yearShort)
}
