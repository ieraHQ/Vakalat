package services

import (
	"context"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/google/uuid"
)

// UserService defines the interface for user operations.
type UserService interface {
	CreateUser(ctx context.Context, user *repositories.User) error
	GetUserByID(ctx context.Context, id string) (*repositories.User, error)
	GetUserByEmail(ctx context.Context, email string) (*repositories.User, error)
	UpdateUser(ctx context.Context, user *repositories.User) error
	DeleteUser(ctx context.Context, id string) error
}

// userService implements UserService.
type userService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(userRepo repositories.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// CreateUser creates a new user.
func (s *userService) CreateUser(ctx context.Context, user *repositories.User) error {
	user.ID = uuid.New().String()
	return s.userRepo.Create(ctx, user)
}

// GetUserByID retrieves a user by ID.
func (s *userService) GetUserByID(ctx context.Context, id string) (*repositories.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// GetUserByEmail retrieves a user by email.
func (s *userService) GetUserByEmail(ctx context.Context, email string) (*repositories.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// UpdateUser updates a user.
func (s *userService) UpdateUser(ctx context.Context, user *repositories.User) error {
	return s.userRepo.Update(ctx, user)
}

// DeleteUser soft-deletes a user.
func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}