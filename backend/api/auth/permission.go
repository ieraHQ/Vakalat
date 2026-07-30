package auth

import (
	"context"
	"errors"

	"github.com/ieraHQ/Vakalat/backend/api/repositories"
)

// PermissionService defines the interface for permission operations.
type PermissionService interface {
	HasPermission(ctx context.Context, userID, permission string) (bool, error)
	HasRole(ctx context.Context, userID, role string) (bool, error)
}

// permissionService implements PermissionService.
type permissionService struct {
	userRepo repositories.UserRepository
}

// NewPermissionService creates a new PermissionService.
func NewPermissionService(userRepo repositories.UserRepository) PermissionService {
	return &permissionService{userRepo: userRepo}
}

// HasPermission checks if a user has a specific permission.
func (s *permissionService) HasPermission(ctx context.Context, userID, permission string) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	// Placeholder: Check user's role permissions (simplified for example)
	if user.RoleID == "11111111-1111-1111-1111-111111111111" { // Admin role
		return true, nil
	}

	return false, nil
}

// HasRole checks if a user has a specific role.
func (s *permissionService) HasRole(ctx context.Context, userID, role string) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	// Placeholder: Check user's role (simplified for example)
	if user.RoleID == role {
		return true, nil
	}

	return false, nil
}