package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ieraHQ/Vakalat/backend/api/auth"
	"github.com/ieraHQ/Vakalat/backend/api/config"
	"github.com/ieraHQ/Vakalat/backend/api/logger"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/ieraHQ/Vakalat/backend/api/services"
	"github.com/ieraHQ/Vakalat/backend/api/validation"
	"go.uber.org/zap"
)

// LoginHandler authenticates a user and issues access/refresh tokens.
func LoginHandler(cfg *config.Config, userService services.UserService, refreshTokenRepo auth.RefreshTokenRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type LoginRequest struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		var req LoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		user, err := userService.GetUserByEmail(c.Context(), req.Email)
		if err != nil {
			logger.GetLogger().Warn("Login failed: user lookup error", zap.String("email", req.Email), zap.Error(err))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		if !auth.VerifyPassword(req.Password, user.PasswordHash) {
			logger.GetLogger().Warn("Login failed: password mismatch", zap.String("email", req.Email), zap.Int("stored_hash_len", len(user.PasswordHash)), zap.Int("submitted_password_len", len(req.Password)))
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		token, err := auth.GenerateToken(user.ID, user.RoleID, cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
		}

		refreshToken, err := auth.GenerateRefreshToken(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate refresh token"})
		}

		if err := refreshTokenRepo.Create(c.Context(), refreshToken); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save refresh token"})
		}

		return c.JSON(fiber.Map{
			"token":         token,
			"refresh_token": refreshToken.Token,
		})
	}
}

// RefreshTokenHandler issues a new access/refresh token pair from a valid refresh token.
func RefreshTokenHandler(cfg *config.Config, refreshTokenRepo auth.RefreshTokenRepository, userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type RefreshRequest struct {
			RefreshToken string `json:"refresh_token"`
		}

		var req RefreshRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		refreshToken, err := refreshTokenRepo.FindByToken(c.Context(), req.RefreshToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid refresh token"})
		}

		if time.Now().After(refreshToken.ExpiresAt) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Refresh token expired"})
		}

		user, err := userService.GetUserByID(c.Context(), refreshToken.UserID)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		token, err := auth.GenerateToken(user.ID, user.RoleID, cfg)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
		}

		newRefreshToken, err := auth.GenerateRefreshToken(user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate refresh token"})
		}

		if err := refreshTokenRepo.Delete(c.Context(), req.RefreshToken); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete refresh token"})
		}

		if err := refreshTokenRepo.Create(c.Context(), newRefreshToken); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save refresh token"})
		}

		return c.JSON(fiber.Map{
			"token":         token,
			"refresh_token": newRefreshToken.Token,
		})
	}
}

// ForgotPasswordHandler issues a password reset token for the given email.
//
// The response never reveals whether the email exists (that would let an
// attacker enumerate accounts), and the raw token is never returned in the
// response — there is no outbound email integration wired up yet, so for now
// the token is only written to the server log. Wire an email/SMS provider
// here before relying on this in production.
func ForgotPasswordHandler(sessionService auth.SessionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type ForgotPasswordRequest struct {
			Email string `json:"email" validate:"required,email"`
		}

		var req ForgotPasswordRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := validation.Validate(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		token, err := sessionService.CreatePasswordResetToken(c.Context(), req.Email)
		if err == nil {
			logger.GetLogger().Info("Password reset token issued (no email provider configured)",
				zap.String("email", req.Email), zap.String("token", token))
		}

		return c.JSON(fiber.Map{"message": "If that email exists, a reset link has been sent"})
	}
}

// ResetPasswordHandler resets a user's password using a reset token.
func ResetPasswordHandler(sessionService auth.SessionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		type ResetPasswordRequest struct {
			Token       string `json:"token" validate:"required"`
			NewPassword string `json:"new_password" validate:"required,min=8"`
		}

		var req ResetPasswordRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := validation.Validate(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		if err := sessionService.ResetPassword(c.Context(), req.Token, req.NewPassword); err != nil {
			if err == auth.ErrInvalidResetToken {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or expired reset token"})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to reset password"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// EnableMFAHandler generates and stores a new TOTP secret for the authenticated user.
func EnableMFAHandler(sessionService auth.SessionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*repositories.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
		}

		secret, otpauthURL, err := sessionService.EnableMFA(c.Context(), user.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to enable MFA"})
		}

		return c.JSON(fiber.Map{"secret": secret, "otpauth_url": otpauthURL})
	}
}

// VerifyMFAHandler verifies a TOTP code for the authenticated user.
func VerifyMFAHandler(sessionService auth.SessionService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*repositories.User)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
		}

		type VerifyMFARequest struct {
			Code string `json:"code" validate:"required,len=6,numeric"`
		}

		var req VerifyMFARequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := validation.Validate(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		valid, err := sessionService.VerifyMFA(c.Context(), user.ID, req.Code)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if !valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid MFA code"})
		}

		return c.JSON(fiber.Map{"verified": true})
	}
}
