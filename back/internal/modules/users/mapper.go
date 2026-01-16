package users

import (
	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

// dbUserToUser converts a database user to domain user
func dbUserToUser(dbUser *db.User) *User {
	var fullName *string
	if dbUser.FullName.Valid {
		fullName = &dbUser.FullName.String
	}

	var cpf *string
	if dbUser.Cpf.Valid {
		cpf = &dbUser.Cpf.String
	}

	// Extract balance
	balanceCents := int64(0)
	if dbUser.BalanceCents.Valid {
		balanceCents = dbUser.BalanceCents.Int64
	}

	// Extract limits
	dailyLimit := int64(100000) // Default R$ 1,000
	if dbUser.DailyTransferLimitCents.Valid {
		dailyLimit = dbUser.DailyTransferLimitCents.Int64
	}

	monthlyLimit := int64(500000) // Default R$ 5,000
	if dbUser.MonthlyTransferLimitCents.Valid {
		monthlyLimit = dbUser.MonthlyTransferLimitCents.Int64
	}

	// Extract status
	status := "active" // Default
	if dbUser.Status.Valid {
		status = dbUser.Status.String
	}

	kycStatus := "pending" // Default
	if dbUser.KycStatus.Valid {
		kycStatus = dbUser.KycStatus.String
	}

	return &User{
		ID:                        dbUser.ID.String(),
		KratosIdentityID:          dbUser.KratosIdentityID,
		Email:                     dbUser.Email,
		FullName:                  fullName,
		CPF:                       cpf,
		BalanceCents:              balanceCents,
		DailyTransferLimitCents:   dailyLimit,
		MonthlyTransferLimitCents: monthlyLimit,
		Status:                    status,
		KYCStatus:                 kycStatus,
		CreatedAt:                 dbUser.CreatedAt.Time,
		UpdatedAt:                 dbUser.UpdatedAt.Time,
	}
}
