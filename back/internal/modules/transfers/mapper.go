package transfers

import (
	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

// dbTransferToTransfer converts a database transfer to domain transfer
func dbTransferToTransfer(dbTransfer *db.Transfer) *Transfer {
	transfer := &Transfer{
		ID:          dbTransfer.ID.String(),
		UserID:      dbTransfer.UserID.String(),
		Type:        dbTransfer.Type,
		Status:      dbTransfer.Status,
		AmountCents: dbTransfer.AmountCents,
		FeeCents:    0,
		Currency:    "BRL",
		CreatedAt:   dbTransfer.CreatedAt.Time,
		UpdatedAt:   dbTransfer.UpdatedAt.Time,
	}

	// Fee
	if dbTransfer.FeeCents.Valid {
		transfer.FeeCents = dbTransfer.FeeCents.Int64
	}

	// Currency
	if dbTransfer.Currency.Valid {
		transfer.Currency = dbTransfer.Currency.String
	}

	// PIX fields
	if dbTransfer.PixKey.Valid {
		key := dbTransfer.PixKey.String
		transfer.PixKey = &key
	}
	if dbTransfer.PixKeyType.Valid {
		keyType := dbTransfer.PixKeyType.String
		transfer.PixKeyType = &keyType
	}

	// TED fields
	if dbTransfer.RecipientName.Valid {
		name := dbTransfer.RecipientName.String
		transfer.RecipientName = &name
	}
	if dbTransfer.RecipientDocument.Valid {
		doc := dbTransfer.RecipientDocument.String
		transfer.RecipientDocument = &doc
	}
	if dbTransfer.RecipientBank.Valid {
		bank := dbTransfer.RecipientBank.String
		transfer.RecipientBank = &bank
	}
	if dbTransfer.RecipientBranch.Valid {
		branch := dbTransfer.RecipientBranch.String
		transfer.RecipientBranch = &branch
	}
	if dbTransfer.RecipientAccount.Valid {
		account := dbTransfer.RecipientAccount.String
		transfer.RecipientAccount = &account
	}
	if dbTransfer.RecipientAccountType.Valid {
		accountType := dbTransfer.RecipientAccountType.String
		transfer.RecipientAccountType = &accountType
	}

	// P2P fields
	if dbTransfer.RecipientUserID.Valid {
		recipientID := dbTransfer.RecipientUserID.UUID.String()
		transfer.RecipientUserID = &recipientID
	}

	// Schedule and completion
	if dbTransfer.ScheduledFor.Valid {
		scheduled := dbTransfer.ScheduledFor.Time
		transfer.ScheduledFor = &scheduled
	}
	if dbTransfer.CompletedAt.Valid {
		completed := dbTransfer.CompletedAt.Time
		transfer.CompletedAt = &completed
	}
	if dbTransfer.FailureReason.Valid {
		reason := dbTransfer.FailureReason.String
		transfer.FailureReason = &reason
	}
	if dbTransfer.AuthenticationCode.Valid {
		code := dbTransfer.AuthenticationCode.String
		transfer.AuthenticationCode = &code
	}

	return transfer
}

// dbTransfersToTransfers converts a slice of database transfers to domain transfers
func dbTransfersToTransfers(dbTransfers []db.Transfer) []Transfer {
	transfers := make([]Transfer, len(dbTransfers))
	for i, dbTransfer := range dbTransfers {
		transfers[i] = *dbTransferToTransfer(&dbTransfer)
	}
	return transfers
}
