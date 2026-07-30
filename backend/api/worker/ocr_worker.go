package worker

import (
	"context"
	"time"

	"github.com/ieraHQ/Vakalat/backend/api/logger"
	"github.com/ieraHQ/Vakalat/backend/api/ocr"
	"github.com/ieraHQ/Vakalat/backend/api/repositories"
	"github.com/ieraHQ/Vakalat/backend/api/services"
	"go.uber.org/zap"
)

// OCRWorker processes OCR jobs in the background.
type OCRWorker struct {
	documentRepo repositories.DocumentRepository
	documentService services.DocumentService
	ocrService ocr.OCRService
	interval time.Duration
}

// NewOCRWorker creates a new OCRWorker.
func NewOCRWorker(
	documentRepo repositories.DocumentRepository,
	documentService services.DocumentService,
	ocrService ocr.OCRService,
	interval time.Duration,
) *OCRWorker {
	return &OCRWorker{
		documentRepo:    documentRepo,
		documentService: documentService,
		ocrService:      ocrService,
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

		logger.GetLogger().Info("OCR completed", zap.String("document_id", doc.ID))
	}
}