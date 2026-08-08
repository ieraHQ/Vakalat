package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ieraHQ/Vakalat/backend/api/services"
)

type summarizeRequest struct {
	MatterID string `json:"matter_id"`
	Text     string `json:"text"`
}

// SummarizeHandler generates and stores an AI summary of the provided text for a matter.
func SummarizeHandler(aiService services.AIService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req summarizeRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if req.MatterID == "" || req.Text == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "matter_id and text are required"})
		}

		summary, err := aiService.Summarize(c.Context(), req.MatterID, req.Text)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate summary"})
		}

		return c.Status(fiber.StatusCreated).JSON(summary)
	}
}

type askRequest struct {
	MatterID string `json:"matter_id"`
	Prompt   string `json:"prompt"`
}

// AskHandler runs a free-form research/Q&A prompt against the AI assistant.
func AskHandler(aiService services.AIService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req askRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if req.Prompt == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "prompt is required"})
		}

		session, err := aiService.Ask(c.Context(), req.MatterID, req.Prompt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get AI response"})
		}

		return c.Status(fiber.StatusCreated).JSON(session)
	}
}

type draftRequest struct {
	MatterID     string `json:"matter_id"`
	Instructions string `json:"instructions"`
}

// DraftHandler generates a draft legal document from instructions.
func DraftHandler(aiService services.AIService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req draftRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		if req.Instructions == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "instructions is required"})
		}

		session, err := aiService.Draft(c.Context(), req.MatterID, req.Instructions)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate draft"})
		}

		return c.Status(fiber.StatusCreated).JSON(session)
	}
}
