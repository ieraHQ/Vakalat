package storage

import (
	"context"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/ieraHQ/Vakalat/backend/api/logger"
)

// Storage defines the interface for file storage operations.
type Storage interface {
	UploadFile(ctx context.Context, file *multipart.FileHeader, matterID string) (string, error)
	DeleteFile(ctx context.Context, filePath string) error
	GetFilePath(documentID string) string
}

// LocalStorage implements Storage for local file storage.
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new LocalStorage.
func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{basePath: basePath}
}

// UploadFile uploads a file to the local filesystem.
func (s *LocalStorage) UploadFile(ctx context.Context, file *multipart.FileHeader, matterID string) (string, error) {
	// Generate a unique filename
	ext := filepath.Ext(file.Filename)
	filename := uuid.New().String() + ext
	filePath := filepath.Join(s.basePath, "documents", matterID, filename)

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	// Open the uploaded file
 src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Copy the file
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	logger.GetLogger().Info("File uploaded", zap.String("path", filePath))
	return filePath, nil
}

// DeleteFile deletes a file from the local filesystem.
func (s *LocalStorage) DeleteFile(ctx context.Context, filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return err
	}

	logger.GetLogger().Info("File deleted", zap.String("path", filePath))
	return nil
}

// GetFilePath returns the full path for a document.
func (s *LocalStorage) GetFilePath(documentID string) string {
	return filepath.Join(s.basePath, "documents", documentID)
}