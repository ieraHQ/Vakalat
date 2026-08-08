package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// UserRepository defines the interface for user operations.
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
	UpdatePasswordHash(ctx context.Context, id, passwordHash string) error
	UpdateMFASecret(ctx context.Context, id, mfaSecret string) error
	GetMFASecret(ctx context.Context, id string) (string, error)
}

// User represents a user in the system.
type User struct {
	ID           string `json:"id"`
	Name         string `json:"name" validate:"required"`
	Email        string `json:"email" validate:"required,email"`
	PasswordHash string `json:"-"`
	RoleID       string `json:"role_id" validate:"required,uuid"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// userRepository implements UserRepository.
type userRepository struct {
	db *pgxpool.Pool
}

// NewUserRepository creates a new UserRepository.
func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

// Create inserts a new user into the database.
func (r *userRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, name, email, password_hash, role_id, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.RoleID, user.IsActive)
	return err
}

// FindByID retrieves a user by ID.
func (r *userRepository) FindByID(ctx context.Context, id string) (*User, error) {
	query := `SELECT id, name, email, password_hash, role_id, is_active, created_at, updated_at FROM users WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.RoleID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByEmail retrieves a user by email.
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, name, email, password_hash, role_id, is_active, created_at, updated_at FROM users WHERE email = $1`
	row := r.db.QueryRow(ctx, query, email)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.RoleID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates a user in the database.
func (r *userRepository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, role_id = $4, is_active = $5, updated_at = NOW()
		WHERE id = $6
	`
	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.PasswordHash, user.RoleID, user.IsActive, user.ID)
	return err
}

// Delete soft-deletes a user from the database.
func (r *userRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE users SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// UpdatePasswordHash updates only a user's password hash.
func (r *userRepository) UpdatePasswordHash(ctx context.Context, id, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, passwordHash, id)
	return err
}

// UpdateMFASecret stores a user's TOTP secret.
func (r *userRepository) UpdateMFASecret(ctx context.Context, id, mfaSecret string) error {
	query := `UPDATE users SET mfa_secret = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, mfaSecret, id)
	return err
}

// GetMFASecret retrieves a user's TOTP secret.
func (r *userRepository) GetMFASecret(ctx context.Context, id string) (string, error) {
	query := `SELECT COALESCE(mfa_secret, '') FROM users WHERE id = $1`
	var secret string
	err := r.db.QueryRow(ctx, query, id).Scan(&secret)
	return secret, err
}