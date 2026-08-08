package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ieraHQ/Vakalat/backend/api/services"
)

// CreateDocumentVersionHandler creates a new version of a document.
func CreateDocumentVersionHandler(documentService services.DocumentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse form data
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form data"})
		}

		files := form.File["file"]
		if len(files) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
		}

		documentID := c.Params("id")
		if documentID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Document ID is required"})
		}

		// Create document version
		version, err := documentService.CreateDocumentVersion(c.Context(), documentID, files[0])
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create document version"})
		}

		return c.Status(fiber.StatusCreated).JSON(version)
	}
}

// ListDocumentVersionsHandler retrieves a list of versions for a document.
func ListDocumentVersionsHandler(documentService services.DocumentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		documentID := c.Params("id")
		if documentID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Document ID is required"})
		}

		versions, err := documentService.ListDocumentVersions(c.Context(), documentID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list document versions"})
		}

		return c.JSON(versions)
	}
}

// GetDocumentVersionHandler retrieves a specific version of a document.
func GetDocumentVersionHandler(documentService services.DocumentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		versionID := c.Params("versionID")
		if versionID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Version ID is required"})
		}

		version, err := documentService.GetDocumentVersion(c.Context(), versionID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document version not found"})
		}

		return c.JSON(version)
	}
}