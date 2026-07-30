package repositories

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DocumentRepository defines the interface for document operations.
type DocumentRepository interface {
	Create(ctx context.Context, document *Document) error
	FindByID(ctx context.Context, id string) (*Document, error)
	Update(ctx context.Context, document *Document) error
	Delete(ctx context.Context, id string) error
	ListByMatter(ctx context.Context, matterID string, limit, offset int) ([]*Document, error)
	UpdateOCRStatus(ctx context.Context, id string, status string) error
	UpdateOCRText(ctx context.Context, id string, text string) error
}

// Document represents a document in the system.
type Document struct {
	ID         string `json:"id"`
	MatterID   string `json:"matter_id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	MimeType   string `json:"mime_type"`
	Size       int64  `json:"size"`
	Hash       string `json:"hash"`
	OCRStatus  string `json:"ocr_status"`
	OCRText    string `json:"ocr_text"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// documentRepository implements DocumentRepository.
type documentRepository struct {
	db *pgxpool.Pool
}

// NewDocumentRepository creates a new DocumentRepository.
func NewDocumentRepository(db *pgxpool.Pool) DocumentRepository {
	return &documentRepository{db: db}
}

// Create inserts a new document into the database.
func (r *documentRepository) Create(ctx context.Context, document *Document) error {
	query := `
		INSERT INTO documents (id, matter_id, name, path, mime_type, size, hash, ocr_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query, document.ID, document.MatterID, document.Name, document.Path, document.MimeType, document.Size, document.Hash, document.OCRStatus)
	return err
}

// FindByID retrieves a document by ID.
func (r *documentRepository) FindByID(ctx context.Context, id string) (*Document, error) {
	query := `SELECT id, matter_id, name, path, mime_type, size, hash, ocr_status, ocr_text, created_at, updated_at FROM documents WHERE id = $1 AND deleted_at IS NULL`
	row := r.db.QueryRow(ctx, query, id)

	var document Document
	err := row.Scan(&document.ID, &document.MatterID, &document.Name, &document.Path, &document.MimeType, &document.Size, &document.Hash, &document.OCRStatus, &document.OCRText, &document.CreatedAt, &document.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &document, nil
}

// Update updates a document in the database.
func (r *documentRepository) Update(ctx context.Context, document *Document) error {
	query := `
		UPDATE documents
		SET name = $1, path = $2, mime_type = $3, size = $4, hash = $5, updated_at = NOW()
		WHERE id = $6
	`
	_, err := r.db.Exec(ctx, query, document.Name, document.Path, document.MimeType, document.Size, document.Hash, document.ID)
	return err
}

// Delete soft-deletes a document from the database.
func (r *documentRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE documents SET deleted_at = NOW() WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// ListByMatter retrieves a list of documents for a matter with pagination.
func (r *documentRepository) ListByMatter(ctx context.Context, matterID string, limit, offset int) ([]*Document, error) {
	query := `SELECT id, matter_id, name, path, mime_type, size, hash, ocr_status, ocr_text, created_at, updated_at FROM documents WHERE matter_id = $1 AND deleted_at IS NULL LIMIT $2 OFFSET $3`
	rows, err := r.db.Query(ctx, query, matterID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []*Document
	for rows.Next() {
		var document Document
		if err := rows.Scan(&document.ID, &document.MatterID, &document.Name, &document.Path, &document.MimeType, &document.Size, &document.Hash, &document.OCRStatus, &document.OCRText, &document.CreatedAt, &document.UpdatedAt); err != nil {
			return nil, err
		}
		documents = append(documents, &document)
	}

	return documents, nil
}

// UpdateOCRStatus updates the OCR status of a document.
func (r *documentRepository) UpdateOCRStatus(ctx context.Context, id string, status string) error {
	query := `UPDATE documents SET ocr_status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, status, id)
	return err
}

// UpdateOCRText updates the OCR text of a document.
func (r *documentRepository) UpdateOCRText(ctx context.Context, id string, text string) error {
	query := `UPDATE documents SET ocr_text = $1, ocr_status = 'completed', updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, text, id)
	return err
}