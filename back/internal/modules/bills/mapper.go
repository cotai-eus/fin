package bills

import (
	"database/sql"
	"time"

	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

// dbBillToBill converts a database bill to domain bill
func dbBillToBill(dbBill *db.Bill) *Bill {
	bill := &Bill{
		ID:               dbBill.ID.String(),
		UserID:           dbBill.UserID.String(),
		Type:             dbBill.Type,
		Status:           dbBill.Status,
		Barcode:          dbBill.Barcode,
		AmountCents:      dbBill.AmountCents,
		RecipientName:    dbBill.RecipientName,
		DueDate:          dbBill.DueDate,
		FinalAmountCents: dbBill.FinalAmountCents,
	}

	if dbBill.CreatedAt.Valid {
		bill.CreatedAt = dbBill.CreatedAt.Time
	}

	if dbBill.FeeCents.Valid {
		bill.FeeCents = dbBill.FeeCents.Int64
	}

	if dbBill.PaymentDate.Valid {
		paymentDate := dbBill.PaymentDate.Time
		bill.PaymentDate = &paymentDate
	}

	return bill
}

// dbBillToBillSummary converts a database bill to bill summary
func dbBillToBillSummary(dbBill *db.Bill) *BillSummary {
	summary := &BillSummary{
		ID:               dbBill.ID.String(),
		Type:             dbBill.Type,
		Status:           dbBill.Status,
		RecipientName:    dbBill.RecipientName,
		AmountCents:      dbBill.AmountCents,
		FinalAmountCents: dbBill.FinalAmountCents,
		DueDate:          dbBill.DueDate,
	}

	if dbBill.CreatedAt.Valid {
		summary.CreatedAt = dbBill.CreatedAt.Time
	}

	return summary
}

// billsToBillSummaries converts multiple database bills to bill summaries
func billsToBillSummaries(dbBills []db.Bill) []*BillSummary {
	summaries := make([]*BillSummary, len(dbBills))
	for i, dbBill := range dbBills {
		summaries[i] = dbBillToBillSummary(&dbBill)
	}
	return summaries
}

// sqlNullTime converts a time pointer to sql.NullTime
func sqlNullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// sqlNullString converts a string pointer to sql.NullString
func sqlNullString(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}
