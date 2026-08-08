package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ieraHQ/Vakalat/backend/api/auth"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/ieraHQ/Vakalat/backend/api/services"
	"github.com/ieraHQ/Vakalat/backend/api/validation"
)

// GetUserHandler retrieves a user by ID.
func GetUserHandler(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		user, err := userService.GetUserByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.JSON(user)
	}
}

// createUserRequest is the client-facing payload for creating a user.
// Password is accepted as plaintext here and hashed before storage —
// repositories.User.PasswordHash is never bound directly from a request body.
type createUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	RoleID   string `json:"role_id" validate:"required,uuid"`
}

// CreateUserHandler creates a new user.
func CreateUserHandler(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req createUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := validation.Validate(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		passwordHash, err := auth.HashPassword(req.Password)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process password"})
		}

		user := repositories.User{
			Name:         req.Name,
			Email:        req.Email,
			PasswordHash: passwordHash,
			RoleID:       req.RoleID,
			IsActive:     true,
		}

		if err := userService.CreateUser(c.Context(), &user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
		}
		return c.Status(fiber.StatusCreated).JSON(user)
	}
}

// updateUserRequest is the client-facing payload for updating a user.
// Password is optional; when present it is hashed before storage.
type updateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"omitempty,min=8"`
	RoleID   string `json:"role_id" validate:"required,uuid"`
	IsActive bool   `json:"is_active"`
}

// UpdateUserHandler updates a user.
func UpdateUserHandler(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var req updateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := validation.Validate(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		user := repositories.User{
			ID:       id,
			Name:     req.Name,
			Email:    req.Email,
			RoleID:   req.RoleID,
			IsActive: req.IsActive,
		}

		if req.Password != "" {
			passwordHash, err := auth.HashPassword(req.Password)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to process password"})
			}
			user.PasswordHash = passwordHash
		} else {
			existing, err := userService.GetUserByID(c.Context(), id)
			if err != nil {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
			}
			user.PasswordHash = existing.PasswordHash
		}

		if err := userService.UpdateUser(c.Context(), &user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
		}
		return c.JSON(user)
	}
}

// DeleteUserHandler deletes a user.
func DeleteUserHandler(userService services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := userService.DeleteUser(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}
