package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Role represents a role in the system.
type Role struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Permission represents a permission in the system.
type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// RoleRepository defines the interface for role and permission lookups.
type RoleRepository interface {
	FindByID(ctx context.Context, id string) (*Role, error)
	FindByName(ctx context.Context, name string) (*Role, error)
	ListPermissionsForRole(ctx context.Context, roleID string) ([]*Permission, error)
}

// roleRepository implements RoleRepository.
type roleRepository struct {
	db *pgxpool.Pool
}

// NewRoleRepository creates a new RoleRepository.
func NewRoleRepository(db *pgxpool.Pool) RoleRepository {
	return &roleRepository{db: db}
}

// FindByID retrieves a role by ID.
func (r *roleRepository) FindByID(ctx context.Context, id string) (*Role, error) {
	query := `SELECT id, name, description FROM roles WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var role Role
	if err := row.Scan(&role.ID, &role.Name, &role.Description); err != nil {
		return nil, err
	}

	return &role, nil
}

// FindByName retrieves a role by name.
func (r *roleRepository) FindByName(ctx context.Context, name string) (*Role, error) {
	query := `SELECT id, name, description FROM roles WHERE name = $1`
	row := r.db.QueryRow(ctx, query, name)

	var role Role
	if err := row.Scan(&role.ID, &role.Name, &role.Description); err != nil {
		return nil, err
	}

	return &role, nil
}

// ListPermissionsForRole retrieves all permissions granted to a role.
func (r *roleRepository) ListPermissionsForRole(ctx context.Context, roleID string) ([]*Permission, error) {
	query := `
		SELECT p.id, p.name, p.description
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`
	rows, err := r.db.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []*Permission
	for rows.Next() {
		var p Permission
		if err := rows.Scan(&p.ID, &p.Name, &p.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, &p)
	}

	return permissions, nil
}
