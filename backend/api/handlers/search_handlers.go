package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ieraHQ/Vakalat/backend/api/ai"
	"github.com/ieraHQ/Vakalat/backend/api/logger"
	"github.com/ieraHQ/Vakalat/backend/api/services"
	"go.uber.org/zap"
)

// SearchHandler runs a hybrid keyword + semantic search over indexed content.
// Query params: q (required), limit (optional, default 10).
func SearchHandler(searchService services.SearchService, embedder ai.LLMClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		query := c.Query("q")
		if query == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Query parameter 'q' is required"})
		}
		limit := c.QueryInt("limit", 10)

		// Semantic (vector) search blends in once the query is embedded. If
		// the embedding model is unreachable, fall back to keyword-only
		// rather than failing the whole search.
		var embedding []float32
		if vec, err := embedder.Embed(c.Context(), query); err != nil {
			logger.GetLogger().Warn("Failed to embed search query, falling back to keyword-only", zap.Error(err))
		} else {
			embedding = vec
		}

		results, err := searchService.Search(c.Context(), query, embedding, limit)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Search failed"})
		}

		return c.JSON(results)
	}
}
