package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/ieraHQ/Vakalat/backend/api/repositories"
)

// fakeUserRepo is a minimal in-memory UserRepository for testing PermissionService
// without a real database.
type fakeUserRepo struct {
	users map[string]*repositories.User
}

func (f *fakeUserRepo) Create(ctx context.Context, user *repositories.User) error { return nil }
func (f *fakeUserRepo) FindByID(ctx context.Context, id string) (*repositories.User, error) {
	u, ok := f.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}
func (f *fakeUserRepo) FindByEmail(ctx context.Context, email string) (*repositories.User, error) {
	return nil, errors.New("not implemented")
}
func (f *fakeUserRepo) Update(ctx context.Context, user *repositories.User) error { return nil }
func (f *fakeUserRepo) Delete(ctx context.Context, id string) error              { return nil }
func (f *fakeUserRepo) UpdatePasswordHash(ctx context.Context, id, passwordHash string) error {
	return nil
}
func (f *fakeUserRepo) UpdateMFASecret(ctx context.Context, id, mfaSecret string) error { return nil }
func (f *fakeUserRepo) GetMFASecret(ctx context.Context, id string) (string, error)     { return "", nil }

// fakeRoleRepo is a minimal in-memory RoleRepository for testing.
type fakeRoleRepo struct {
	roles       map[string]*repositories.Role
	permissions map[string][]*repositories.Permission // roleID -> permissions
}

func (f *fakeRoleRepo) FindByID(ctx context.Context, id string) (*repositories.Role, error) {
	r, ok := f.roles[id]
	if !ok {
		return nil, errors.New("role not found")
	}
	return r, nil
}
func (f *fakeRoleRepo) FindByName(ctx context.Context, name string) (*repositories.Role, error) {
	for _, r := range f.roles {
		if r.Name == name {
			return r, nil
		}
	}
	return nil, errors.New("role not found")
}
func (f *fakeRoleRepo) ListPermissionsForRole(ctx context.Context, roleID string) ([]*repositories.Permission, error) {
	return f.permissions[roleID], nil
}

func newTestPermissionService() PermissionService {
	userRepo := &fakeUserRepo{
		users: map[string]*repositories.User{
			"user-admin": {ID: "user-admin", RoleID: "role-admin"},
			"user-staff": {ID: "user-staff", RoleID: "role-staff"},
			"user-none":  {ID: "user-none", RoleID: ""},
		},
	}
	roleRepo := &fakeRoleRepo{
		roles: map[string]*repositories.Role{
			"role-admin": {ID: "role-admin", Name: "admin"},
			"role-staff": {ID: "role-staff", Name: "staff"},
		},
		permissions: map[string][]*repositories.Permission{
			"role-admin": {{ID: "p1", Name: "manage_users"}, {ID: "p2", Name: "manage_matters"}},
			"role-staff": {{ID: "p2", Name: "manage_matters"}},
		},
	}
	return NewPermissionService(userRepo, roleRepo)
}

func TestHasPermission_GrantedByRole(t *testing.T) {
	s := newTestPermissionService()

	ok, err := s.HasPermission(context.Background(), "user-admin", "manage_users")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected admin to have manage_users permission")
	}
}

func TestHasPermission_DeniedWhenRoleLacksIt(t *testing.T) {
	s := newTestPermissionService()

	ok, err := s.HasPermission(context.Background(), "user-staff", "manage_users")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected staff to be denied manage_users permission")
	}
}

func TestHasPermission_DeniedWhenNoRole(t *testing.T) {
	// Regression test: a user with no role must never silently pass a
	// permission check.
	s := newTestPermissionService()

	ok, err := s.HasPermission(context.Background(), "user-none", "manage_users")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected user with no role to be denied any permission")
	}
}

func TestHasRole_MatchesByName(t *testing.T) {
	s := newTestPermissionService()

	ok, err := s.HasRole(context.Background(), "user-admin", "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected user-admin to have role 'admin'")
	}

	// Regression test: HasRole previously compared a role UUID against the
	// role-name string directly, so this could never have matched.
	ok, err = s.HasRole(context.Background(), "user-admin", "role-admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected HasRole to compare by role name, not role ID")
	}
}
