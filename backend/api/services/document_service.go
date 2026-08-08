package services

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/ieraHQ/Vakalat/backend/api/storage"
	"github.com/google/uuid"
)

// DocumentService defines the interface for document operations.
type DocumentService interface {
	CreateDocument(ctx context.Context, matterID string, file *multipart.FileHeader) (*repositories.Document, error)
	GetDocumentByID(ctx context.Context, id string) (*repositories.Document, error)
	UpdateDocument(ctx context.Context, document *repositories.Document) error
	DeleteDocument(ctx context.Context, id string) error
	ListDocumentsByMatter(ctx context.Context, matterID string, limit, offset int) ([]*repositories.Document, error)
	UpdateDocumentOCRStatus(ctx context.Context, id string, status string) error
	UpdateDocumentOCRText(ctx context.Context, id string, text string) error
	CreateDocumentVersion(ctx context.Context, documentID string, file *multipart.FileHeader) (*repositories.DocumentVersion, error)
	ListDocumentVersions(ctx context.Context, documentID string) ([]*repositories.DocumentVersion, error)
	GetDocumentVersion(ctx context.Context, versionID string) (*repositories.DocumentVersion, error)
}

// documentService implements DocumentService.
type documentService struct {
	documentRepo repositories.DocumentRepository
	storage     storage.Storage
}

// NewDocumentService creates a new DocumentService.
func NewDocumentService(documentRepo repositories.DocumentRepository, storage storage.Storage) DocumentService {
	return &documentService{documentRepo: documentRepo, storage: storage}
}

// CreateDocument creates a new document and uploads the file.
func (s *documentService) CreateDocument(ctx context.Context, matterID string, file *multipart.FileHeader) (*repositories.Document, error) {
	// Upload file
	filePath, err := s.storage.UploadFile(ctx, file, matterID)
	if err != nil {
		return nil, err
	}

	// Create document record
	document := &repositories.Document{
		ID:        uuid.New().String(),
		MatterID:  matterID,
		Name:      file.Filename,
		Path:      filePath,
		MimeType:  file.Header.Get("Content-Type"),
		Size:      file.Size,
		OCRStatus: "pending",
		CreatedAt: time.Now(),
	}

	if err := s.documentRepo.Create(ctx, document); err != nil {
		return nil, err
	}

	return document, nil
}

// GetDocumentByID retrieves a document by ID.
func (s *documentService) GetDocumentByID(ctx context.Context, id string) (*repositories.Document, error) {
	return s.documentRepo.FindByID(ctx, id)
}

// UpdateDocument updates a document.
func (s *documentService) UpdateDocument(ctx context.Context, document *repositories.Document) error {
	return s.documentRepo.Update(ctx, document)
}

// DeleteDocument soft-deletes a document and removes the file.
func (s *documentService) DeleteDocument(ctx context.Context, id string) error {
	document, err := s.documentRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := s.storage.DeleteFile(ctx, document.Path); err != nil {
		return err
	}

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

// CreateDocumentVersion creates a new version of a document.
func (s *documentService) CreateDocumentVersion(ctx context.Context, documentID string, file *multipart.FileHeader) (*repositories.DocumentVersion, error) {
	// Get the latest version
	versions, err := s.documentRepo.ListVersions(ctx, documentID)
	if err != nil {
		return nil, err
	}

	latestVersion := 0
	if len(versions) > 0 {
		latestVersion = versions[0].Version
	}

	// Upload new version
	versionPath, err := s.storage.UploadFile(ctx, file, documentID)
	if err != nil {
		return nil, err
	}

	// Create version record
	version := &repositories.DocumentVersion{
		ID:        uuid.New().String(),
		DocumentID: documentID,
		Version:    latestVersion + 1,
		Path:      versionPath,
	}

	if err := s.documentRepo.CreateVersion(ctx, version); err != nil {
		return nil, err
	}

	return version, nil
}

// ListDocumentVersions retrieves a list of versions for a document.
func (s *documentService) ListDocumentVersions(ctx context.Context, documentID string) ([]*repositories.DocumentVersion, error) {
	return s.documentRepo.ListVersions(ctx, documentID)
}

// GetDocumentVersion retrieves a specific version of a document.
func (s *documentService) GetDocumentVersion(ctx context.Context, versionID string) (*repositories.DocumentVersion, error) {
	return s.documentRepo.GetVersion(ctx, versionID)
}