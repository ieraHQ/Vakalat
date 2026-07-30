package services

import (
	"context"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/google/uuid"
)

// DocumentService defines the interface for document operations.
type DocumentService interface {
	CreateDocument(ctx context.Context, document *repositories.Document) error
	GetDocumentByID(ctx context.Context, id string) (*repositories.Document, error)
	UpdateDocument(ctx context.Context, document *repositories.Document) error
	DeleteDocument(ctx context.Context, id string) error
	ListDocumentsByMatter(ctx context.Context, matterID string, limit, offset int) ([]*repositories.Document, error)
	UpdateDocumentOCRStatus(ctx context.Context, id string, status string) error
	UpdateDocumentOCRText(ctx context.Context, id string, text string) error
}

// documentService implements DocumentService.
type documentService struct {
	documentRepo repositories.DocumentRepository
}

// NewDocumentService creates a new DocumentService.
func NewDocumentService(documentRepo repositories.DocumentRepository) DocumentService {
	return &documentService{documentRepo: documentRepo}
}

// CreateDocument creates a new document.
func (s *documentService) CreateDocument(ctx context.Context, document *repositories.Document) error {
	document.ID = uuid.New().String()
	return s.documentRepo.Create(ctx, document)
}

// GetDocumentByID retrieves a document by ID.
func (s *documentService) GetDocumentByID(ctx context.Context, id string) (*repositories.Document, error) {
	return s.documentRepo.FindByID(ctx, id)
}

// UpdateDocument updates a document.
func (s *documentService) UpdateDocument(ctx context.Context, document *repositories.Document) error {
	return s.documentRepo.Update(ctx, document)
}

// DeleteDocument soft-deletes a document.
func (s *documentService) DeleteDocument(ctx context.Context, id string) error {
	return s.documentRepo.Delete(ctx, id)
}

// ListDocumentsByMatter retrieves a list of documents for a matter with pagination.
func (s *documentService) ListDocumentsByMatter(ctx context.Context, matterID string, limit, offset int) ([]*repositories.Document, error) {
	return s.documentRepo.ListByMatter(ctx, matterID, limit, offset)
}

// UpdateDocumentOCRStatus updates the OCR status of a document.
func (s *documentService) UpdateDocumentOCRStatus(ctx context.Context, id string, status string) error {
	return s.documentRepo.UpdateOCRStatus(ctx, id, status)
}

// UpdateDocumentOCRText updates the OCR text of a document.
func (s *documentService) UpdateDocumentOCRText(ctx context.Context, id string, text string) error {
	return s.documentRepo.UpdateOCRText(ctx, id, text)
}