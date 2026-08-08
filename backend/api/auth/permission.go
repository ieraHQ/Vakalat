package auth

import (
	"context"

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
	roleRepo repositories.RoleRepository
}

// NewPermissionService creates a new PermissionService.
func NewPermissionService(userRepo repositories.UserRepository, roleRepo repositories.RoleRepository) PermissionService {
	return &permissionService{userRepo: userRepo, roleRepo: roleRepo}
}

// HasPermission checks whether the user's role grants the given permission.
func (s *permissionService) HasPermission(ctx context.Context, userID, permission string) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	if user.RoleID == "" {
		return false, nil
	}

	permissions, err := s.roleRepo.ListPermissionsForRole(ctx, user.RoleID)
	if err != nil {
		return false, err
	}

	for _, p := range permissions {
		if p.Name == permission {
			return true, nil
		}
	}

	return false, nil
}

// HasRole checks whether the user is assigned the given role (by name).
func (s *permissionService) HasRole(ctx context.Context, userID, role string) (bool, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return false, err
	}

	if user.RoleID == "" {
		return false, nil
	}

	userRole, err := s.roleRepo.FindByID(ctx, user.RoleID)
	if err != nil {
		return false, err
	}

	return userRole.Name == role, nil
}
