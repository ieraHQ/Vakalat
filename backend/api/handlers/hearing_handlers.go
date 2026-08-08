package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/ieraHQ/Vakalat/backend/api/services"
	"github.com/ieraHQ/Vakalat/backend/api/validation"
)

// GetHearingHandler retrieves a hearing by ID.
func GetHearingHandler(hearingService services.HearingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		hearing, err := hearingService.GetHearingByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Hearing not found"})
		}
		return c.JSON(hearing)
	}
}

// CreateHearingHandler creates a new hearing.
func CreateHearingHandler(hearingService services.HearingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var hearing repositories.Hearing
		if err := c.BodyParser(&hearing); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := validation.Validate(&hearing); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := hearingService.CreateHearing(c.Context(), &hearing); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create hearing"})
		}
		return c.Status(fiber.StatusCreated).JSON(hearing)
	}
}

// UpdateHearingHandler updates a hearing.
func UpdateHearingHandler(hearingService services.HearingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var hearing repositories.Hearing
		if err := c.BodyParser(&hearing); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		hearing.ID = id
		if err := validation.Validate(&hearing); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := hearingService.UpdateHearing(c.Context(), &hearing); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update hearing"})
		}
		return c.JSON(hearing)
	}
}

// DeleteHearingHandler deletes a hearing.
func DeleteHearingHandler(hearingService services.HearingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := hearingService.DeleteHearing(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete hearing"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// ListHearingsByMatterHandler retrieves all hearings for a matter.
func ListHearingsByMatterHandler(hearingService services.HearingService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		matterID := c.Params("matterID")
		hearings, err := hearingService.ListHearingsByMatter(c.Context(), matterID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list hearings"})
		}
		return c.JSON(hearings)
	}
}
