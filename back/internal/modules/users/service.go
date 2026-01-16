package users

import (
	"context"
	"database/sql"
	"errors"
)

// Service handles business logic for users
type Service struct {
	repo *Repository
}

// NewService creates a new user service
func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// GetCurrentUser retrieves the current authenticated user
func (s *Service) GetCurrentUser(ctx context.Context, kratosID string) (*User, error) {
	dbUser, err := s.repo.GetByKratosID(ctx, kratosID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return dbUserToUser(dbUser), nil
}

// GetByID retrieves a user by ID
func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	dbUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return dbUserToUser(dbUser), nil
}

// Create creates a new user
func (s *Service) Create(ctx context.Context, req CreateUserRequest) (*User, error) {
	// Check if user already exists
	existingUser, err := s.repo.GetByKratosID(ctx, req.KratosIdentityID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	dbUser, err := s.repo.Create(ctx, req)
	if err != nil {
		return nil, err
	}

	return dbUserToUser(dbUser), nil
}

// Update updates a user
func (s *Service) Update(ctx context.Context, id string, req UpdateUserRequest) (*User, error) {
	// Verify user exists
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	dbUser, err := s.repo.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return dbUserToUser(dbUser), nil
}
