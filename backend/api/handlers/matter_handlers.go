package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/ieraHQ/Vakalat/backend/api/services"
	"github.com/ieraHQ/Vakalat/backend/api/validation"
)

// GetMatterHandler retrieves a matter by ID.
func GetMatterHandler(matterService services.MatterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		matter, err := matterService.GetMatterByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Matter not found"})
		}
		return c.JSON(matter)
	}
}

// CreateMatterHandler creates a new matter.
func CreateMatterHandler(matterService services.MatterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var matter repositories.Matter
		if err := c.BodyParser(&matter); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if err := validation.Validate(&matter); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := matterService.CreateMatter(c.Context(), &matter); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create matter"})
		}
		return c.Status(fiber.StatusCreated).JSON(matter)
	}
}

// UpdateMatterHandler updates a matter.
func UpdateMatterHandler(matterService services.MatterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var matter repositories.Matter
		if err := c.BodyParser(&matter); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		matter.ID = id
		if err := validation.Validate(&matter); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := matterService.UpdateMatter(c.Context(), &matter); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update matter"})
		}
		return c.JSON(matter)
	}
}

// DeleteMatterHandler deletes a matter.
func DeleteMatterHandler(matterService services.MatterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := matterService.DeleteMatter(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete matter"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// ListMattersHandler retrieves all matters with pagination.
func ListMattersHandler(matterService services.MatterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		limit := c.QueryInt("limit", 20)
		offset := c.QueryInt("offset", 0)
		matters, err := matterService.ListMatters(c.Context(), limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list matters"})
		}
		return c.JSON(matters)
	}
}

// ListMattersByClientHandler retrieves a list of matters for a client.
func ListMattersByClientHandler(matterService services.MatterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		clientID := c.Params("clientID")
		limit := c.QueryInt("limit", 10)
		offset := c.QueryInt("offset", 0)
		matters, err := matterService.ListMattersByClient(c.Context(), clientID, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list matters"})
		}
		return c.JSON(matters)
	}
}

// ListMattersByAdvocateHandler retrieves a list of matters for an advocate.
func ListMattersByAdvocateHandler(matterService services.MatterService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		advocateID := c.Params("advocateID")
		limit := c.QueryInt("limit", 10)
		offset := c.QueryInt("offset", 0)
		matters, err := matterService.ListMattersByAdvocate(c.Context(), advocateID, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list matters"})
		}
		return c.JSON(matters)
	}
}
