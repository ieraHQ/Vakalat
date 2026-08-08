package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/pquerna/otp/totp"
)

const passwordResetTokenTTL = 1 * time.Hour

var (
	// ErrInvalidResetToken is returned when a reset token is missing, expired, or already used.
	ErrInvalidResetToken = errors.New("invalid or expired reset token")
	// ErrMFANotEnabled is returned when verifying a code for a user with no MFA secret set.
	ErrMFANotEnabled = errors.New("mfa is not enabled for this user")
)

// SessionService defines the interface for session operations.
type SessionService interface {
	CreatePasswordResetToken(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token, newPassword string) error
	EnableMFA(ctx context.Context, userID string) (secret, otpauthURL string, err error)
	VerifyMFA(ctx context.Context, userID, code string) (bool, error)
}

// sessionService implements SessionService.
type sessionService struct {
	userRepo          repositories.UserRepository
	passwordResetRepo repositories.PasswordResetRepository
}

// NewSessionService creates a new SessionService.
func NewSessionService(userRepo repositories.UserRepository, passwordResetRepo repositories.PasswordResetRepository) SessionService {
	return &sessionService{userRepo: userRepo, passwordResetRepo: passwordResetRepo}
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// CreatePasswordResetToken issues a one-time password reset token for the given email.
// The raw token is returned only to the caller (e.g. to email to the user) — only its
// hash is persisted, so a leaked database can't be used to reset passwords directly.
func (s *sessionService) CreatePasswordResetToken(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	token := uuid.New().String()
	resetToken := &repositories.PasswordResetToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: hashToken(token),
		ExpiresAt: time.Now().Add(passwordResetTokenTTL),
	}

	if err := s.passwordResetRepo.Create(ctx, resetToken); err != nil {
		return "", err
	}

	return token, nil
}

// ResetPassword validates a reset token and, if valid and unused, updates the user's password.
func (s *sessionService) ResetPassword(ctx context.Context, token, newPassword string) error {
	resetToken, err := s.passwordResetRepo.FindByTokenHash(ctx, hashToken(token))
	if err != nil {
		return ErrInvalidResetToken
	}

	if resetToken.UsedAt != nil || time.Now().After(resetToken.ExpiresAt) {
		return ErrInvalidResetToken
	}

	passwordHash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	if err := s.userRepo.UpdatePasswordHash(ctx, resetToken.UserID, passwordHash); err != nil {
		return err
	}

	return s.passwordResetRepo.MarkUsed(ctx, resetToken.ID)
}

// EnableMFA generates a new TOTP secret for the user, stores it, and returns the
// secret plus an otpauth:// URL the client can render as a QR code.
func (s *sessionService) EnableMFA(ctx context.Context, userID string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "Vakalat",
		AccountName: userID,
	})
	if err != nil {
		return "", "", err
	}

	if err := s.userRepo.UpdateMFASecret(ctx, userID, key.Secret()); err != nil {
		return "", "", err
	}

	return key.Secret(), key.URL(), nil
}

// VerifyMFA checks a TOTP code against the user's stored secret.
func (s *sessionService) VerifyMFA(ctx context.Context, userID, code string) (bool, error) {
	secret, err := s.userRepo.GetMFASecret(ctx, userID)
	if err != nil {
		return false, err
	}
	if secret == "" {
		return false, ErrMFANotEnabled
	}

	return totp.Validate(code, secret), nil
}
