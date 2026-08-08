package ocr

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ieraHQ/Vakalat/backend/api/logger"
	"go.uber.org/zap"
)

// OCRService defines the interface for OCR operations.
type OCRService interface {
	ExtractText(ctx context.Context, filePath string) (string, error)
}

// ocrService implements OCRService.
type ocrService struct {
	useOCRmyPDF bool
	usePaddleOCR bool
}

// NewOCRService creates a new OCRService.
func NewOCRService(useOCRmyPDF, usePaddleOCR bool) OCRService {
	return &ocrService{
		useOCRmyPDF: useOCRmyPDF,
		usePaddleOCR: usePaddleOCR,
	}
}

// ExtractText extracts text from a file using OCR.
func (s *ocrService) ExtractText(ctx context.Context, filePath string) (string, error) {
	// Try OCRmyPDF first (best for PDFs)
	if s.useOCRmyPDF && strings.ToLower(filepath.Ext(filePath)) == ".pdf" {
		text, err := s.extractTextWithOCRmyPDF(filePath)
		if err == nil {
			return text, nil
		}
		logger.GetLogger().Warn("OCRmyPDF failed, falling back", zap.Error(err))
	}

	// Try PaddleOCR (best for images)
	if s.usePaddleOCR {
		text, err := s.extractTextWithPaddleOCR(filePath)
		if err == nil {
			return text, nil
		}
		logger.GetLogger().Warn("PaddleOCR failed, falling back", zap.Error(err))
	}

	// Fallback to Tesseract
	return s.extractTextWithTesseract(filePath)
}

// extractTextWithOCRmyPDF extracts text from a PDF using OCRmyPDF.
func (s *ocrService) extractTextWithOCRmyPDF(filePath string) (string, error) {
	// OCRmyPDF outputs a text file
	outputFile := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".txt"

	// Run OCRmyPDF command
	cmd := exec.Command(
		"ocrmypdf",
		"--force-ocr",
		"--sidecar", outputFile,
		filePath,
		"-", // Output to stdout (discarded)
	)

	// Execute command
	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Read the extracted text
	text, err := os.ReadFile(outputFile)
	if err != nil {
		return "", err
	}

	// Clean up
	os.Remove(outputFile)

	return string(text), nil
}

// extractTextWithPaddleOCR extracts text from an image using PaddleOCR.
func (s *ocrService) extractTextWithPaddleOCR(filePath string) (string, error) {
	// Run PaddleOCR command
	cmd := exec.Command("paddleocr", "--image_dir", filePath, "--use_angle_cls", "true")

	// Execute command
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// extractTextWithTesseract extracts text using Tesseract (fallback).
func (s *ocrService) extractTextWithTesseract(filePath string) (string, error) {
	// Run Tesseract command
	cmd := exec.Command("tesseract", filePath, "stdout")

	// Execute command
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}