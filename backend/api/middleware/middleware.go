package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ieraHQ/Vakalat/backend/api/logger"
	"github.com/ieraHQ/Vakalat/backend/api/services"
)

// AuthMiddleware validates JWT tokens and attaches the user to the context.
func AuthMiddleware(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header missing"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Bearer token missing"})
		}

		// Validate token (placeholder for JWT validation)
		user, err := userService.GetUserByID(c.Context(), token) // Simplified for example
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		// Attach user to context
		c.Locals("user", user)
		return c.Next()
	}
}

// RBACMiddleware checks if the user has the required permission.
func RBACMiddleware(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*services.User)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
		}

		// Check permission (placeholder for RBAC logic)
		if user.Role != "admin" { // Simplified for example
			return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Permission denied"})
		}

		return c.Next()
	}
}

// LoggerMiddleware logs HTTP requests.
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.GetLogger().Info(
			"HTTP Request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
		)
		return c.Next()
	}
}