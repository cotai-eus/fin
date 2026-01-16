package users

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	db "github.com/lauratech/fin/back/internal/shared/database/sqlc"
)

// Repository handles data access for users
type Repository struct {
	queries *db.Queries
}

// NewRepository creates a new user repository
func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		queries: db.New(database),
	}
}

// GetByID retrieves a user by ID
func (r *Repository) GetByID(ctx context.Context, id string) (*db.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	user, err := r.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByKratosID retrieves a user by Kratos identity ID
func (r *Repository) GetByKratosID(ctx context.Context, kratosID string) (*db.User, error) {
	user, err := r.queries.GetUserByKratosID(ctx, kratosID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *Repository) GetByEmail(ctx context.Context, email string) (*db.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Create creates a new user
func (r *Repository) Create(ctx context.Context, req CreateUserRequest) (*db.User, error) {
	// Convert optional fields to sql.NullString
	fullName := sql.NullString{Valid: false}
	if req.FullName != nil {
		fullName = sql.NullString{String: *req.FullName, Valid: true}
	}

	cpf := sql.NullString{Valid: false}
	if req.CPF != nil {
		cpf = sql.NullString{String: *req.CPF, Valid: true}
	}

	user, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		KratosIdentityID: req.KratosIdentityID,
		Email:            req.Email,
		FullName:         fullName,
		Cpf:              cpf,
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates a user
func (r *Repository) Update(ctx context.Context, id string, req UpdateUserRequest) (*db.User, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	fullName := sql.NullString{Valid: false}
	if req.FullName != nil {
		fullName = sql.NullString{String: *req.FullName, Valid: true}
	}

	user, err := r.queries.UpdateUser(ctx, db.UpdateUserParams{
		ID:       userID,
		FullName: fullName,
	})
	if err != nil {
		return nil, err
	}

	return &user, nil
}
