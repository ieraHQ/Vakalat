package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PasswordResetToken represents a password reset token in the system.
type PasswordResetToken struct {
	ID        string
	UserID    string
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
}

// PasswordResetRepository defines the interface for password reset token persistence.
type PasswordResetRepository interface {
	Create(ctx context.Context, t *PasswordResetToken) error
	FindByTokenHash(ctx context.Context, tokenHash string) (*PasswordResetToken, error)
	MarkUsed(ctx context.Context, id string) error
}

// passwordResetRepository implements PasswordResetRepository.
type passwordResetRepository struct {
	db *pgxpool.Pool
}

// NewPasswordResetRepository creates a new PasswordResetRepository.
func NewPasswordResetRepository(db *pgxpool.Pool) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

// Create inserts a new password reset token.
func (r *passwordResetRepository) Create(ctx context.Context, t *PasswordResetToken) error {
	query := `INSERT INTO password_reset_tokens (id, user_id, token_hash, expires_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, t.ID, t.UserID, t.TokenHash, t.ExpiresAt)
	return err
}

// FindByTokenHash retrieves a password reset token by its hashed value.
func (r *passwordResetRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*PasswordResetToken, error) {
	query := `SELECT id, user_id, token_hash, expires_at, used_at FROM password_reset_tokens WHERE token_hash = $1`
	row := r.db.QueryRow(ctx, query, tokenHash)

	var t PasswordResetToken
	if err := row.Scan(&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.UsedAt); err != nil {
		return nil, err
	}
	return &t, nil
}

// MarkUsed marks a password reset token as consumed.
func (r *passwordResetRepository) MarkUsed(ctx context.Context, id string) error {
	query := `UPDATE password_reset_tokens SET used_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
