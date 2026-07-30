package middleware

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ieraHQ/Vakalat/backend/api/auth"
	"github.com/ieraHQ/Vakalat/backend/api/config"
	"github.com/ieraHQ/Vakalat/backend/api/logger"
	"github.com/ieraHQ/Vakalat/backend/api/services"
)

// AuthMiddleware validates JWT tokens and attaches the user to the context.
func AuthMiddleware(cfg *config.Config, userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header missing"})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Bearer token missing"})
		}

		claims, err := auth.ValidateToken(token, cfg)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		user, err := userService.GetUserByID(c.Context(), claims.UserID)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
		}

		// Attach user to context
		c.Locals("user", user)
		return c.Next()
	}
}

// RBACMiddleware checks if the user has the required permission.
func RBACMiddleware(permissionService auth.PermissionService, permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(*services.User)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
		}

		hasPermission, err := permissionService.HasPermission(c.Context(), user.ID, permission)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check permission"})
		}

		if !hasPermission {
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