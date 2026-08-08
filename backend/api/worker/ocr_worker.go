package worker

import (
	"context"
	"time"

	"github.com/ieraHQ/Vakalat/backend/api/ai"
	"github.com/ieraHQ/Vakalat/backend/api/logger"
	"github.com/ieraHQ/Vakalat/backend/api/ocr"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/ieraHQ/Vakalat/backend/api/services"
	"go.uber.org/zap"
)

// OCRWorker processes OCR jobs in the background.
type OCRWorker struct {
	documentRepo    repositories.DocumentRepository
	documentService services.DocumentService
	ocrService      ocr.OCRService
	searchService   services.SearchService
	embedder        ai.LLMClient
	interval        time.Duration
}

// NewOCRWorker creates a new OCRWorker. searchService and embedder index a
// document's extracted text into the search index once OCR completes, so
// document content becomes searchable without a separate backfill step.
func NewOCRWorker(
	documentRepo repositories.DocumentRepository,
	documentService services.DocumentService,
	ocrService ocr.OCRService,
	searchService services.SearchService,
	embedder ai.LLMClient,
	interval time.Duration,
) *OCRWorker {
	return &OCRWorker{
		documentRepo:    documentRepo,
		documentService: documentService,
		ocrService:      ocrService,
		searchService:   searchService,
		embedder:        embedder,
		interval:        interval,
	}
}

// Start begins processing OCR jobs.
func (w *OCRWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.processPendingDocuments(ctx)
		}
	}
}

// processPendingDocuments processes documents with OCR status "pending".
func (w *OCRWorker) processPendingDocuments(ctx context.Context) {
	documents, err := w.documentRepo.ListByOCRStatus(ctx, "pending")
	if err != nil {
		logger.GetLogger().Error("Failed to list pending documents", zap.Error(err))
		return
	}

	for _, doc := range documents {
		// Update status to "processing"
		if err := w.documentService.UpdateDocumentOCRStatus(ctx, doc.ID, "processing"); err != nil {
			logger.GetLogger().Error("Failed to update OCR status", zap.String("document_id", doc.ID), zap.Error(err))
			continue
		}

		// Extract text
		text, err := w.ocrService.ExtractText(ctx, doc.Path)
		if err != nil {
			logger.GetLogger().Error("Failed to extract text", zap.String("document_id", doc.ID), zap.Error(err))
			if err := w.documentService.UpdateDocumentOCRStatus(ctx, doc.ID, "failed"); err != nil {
				logger.GetLogger().Error("Failed to update OCR status", zap.String("document_id", doc.ID), zap.Error(err))
			}
			continue
		}

		// Update document with extracted text
		if err := w.documentService.UpdateDocumentOCRText(ctx, doc.ID, text); err != nil {
			logger.GetLogger().Error("Failed to update OCR text", zap.String("document_id", doc.ID), zap.Error(err))
			continue
		}

		w.indexDocument(ctx, doc.ID, doc.Name, text)

		logger.GetLogger().Info("OCR completed", zap.String("document_id", doc.ID))
	}
}

// indexDocument embeds and upserts a document's OCR text into the search
// index. Failures are logged, not fatal — a document that fails to index can
// still be retried on the next backfill run.
func (w *OCRWorker) indexDocument(ctx context.Context, documentID, name, text string) {
	content := name + "\n" + text

	embedding, err := w.embedder.Embed(ctx, content)
	if err != nil {
		logger.GetLogger().Warn("Failed to embed document for search index", zap.String("document_id", documentID), zap.Error(err))
		embedding = nil
	}

	if err := w.searchService.IndexContent(ctx, "document", documentID, content, embedding); err != nil {
		logger.GetLogger().Warn("Failed to index document", zap.String("document_id", documentID), zap.Error(err))
	}
}
