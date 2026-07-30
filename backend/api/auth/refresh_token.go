package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
)

// RefreshTokenRepository defines the interface for refresh token operations.
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *RefreshToken) error
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)
	Delete(ctx context.Context, token string) error
}

// RefreshToken represents a refresh token in the system.
type RefreshToken struct {
	Token     string `json:"token"`
	UserID    string `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

// refreshTokenRepository implements RefreshTokenRepository.
type refreshTokenRepository struct {
	db repositories.DB
}

// NewRefreshTokenRepository creates a new RefreshTokenRepository.
func NewRefreshTokenRepository(db repositories.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

// Create inserts a new refresh token into the database.
func (r *refreshTokenRepository) Create(ctx context.Context, token *RefreshToken) error {
	query := `INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, token.Token, token.UserID, token.ExpiresAt)
	return err
}

// FindByToken retrieves a refresh token by token string.
func (r *refreshTokenRepository) FindByToken(ctx context.Context, token string) (*RefreshToken, error) {
	query := `SELECT token, user_id, expires_at FROM refresh_tokens WHERE token = $1`
	row := r.db.QueryRow(ctx, query, token)

	var refreshToken RefreshToken
	err := row.Scan(&refreshToken.Token, &refreshToken.UserID, &refreshToken.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &refreshToken, nil
}

// Delete removes a refresh token from the database.
func (r *refreshTokenRepository) Delete(ctx context.Context, token string) error {
	query := `DELETE FROM refresh_tokens WHERE token = $1`
	_, err := r.db.Exec(ctx, query, token)
	return err
}

// GenerateRefreshToken generates a new refresh token for the given user.
func GenerateRefreshToken(userID string) (*RefreshToken, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 7 days

	return &RefreshToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: expiresAt,
	}, nil
}