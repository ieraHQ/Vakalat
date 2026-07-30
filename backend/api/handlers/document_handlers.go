package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ieraHQ/Vakalat/backend/api/services"
)

// GetDocumentHandler retrieves a document by ID.
func GetDocumentHandler(documentService services.DocumentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		document, err := documentService.GetDocumentByID(c.Context(), id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Document not found"})
		}
		return c.JSON(document)
	}
}

// CreateDocumentHandler creates a new document.
func CreateDocumentHandler(documentService services.DocumentService) fiber.Handler {
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

		matterID := c.FormValue("matter_id")
		if matterID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Matter ID is required"})
		}

		// Create document
		document, err := documentService.CreateDocument(c.Context(), matterID, files[0])
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create document"})
		}

		return c.Status(fiber.StatusCreated).JSON(document)
	}
}

// UpdateDocumentHandler updates a document.
func UpdateDocumentHandler(documentService services.DocumentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var document repositories.Document
		if err := c.BodyParser(&document); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}
		document.ID = id
		if err := documentService.UpdateDocument(c.Context(), &document); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update document"})
		}
		return c.JSON(document)
	}
}

// DeleteDocumentHandler deletes a document.
func DeleteDocumentHandler(documentService services.DocumentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		if err := documentService.DeleteDocument(c.Context(), id); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete document"})
		}
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// ListDocumentsByMatterHandler retrieves a list of documents for a matter.
func ListDocumentsByMatterHandler(documentService services.DocumentService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		matterID := c.Params("matterID")
		limit := c.QueryInt("limit", 10)
		offset := c.QueryInt("offset", 0)
		documents, err := documentService.ListDocumentsByMatter(c.Context(), matterID, limit, offset)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list documents"})
		}
		return c.JSON(documents)
	}
}