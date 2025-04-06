package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type LocalStorage struct {
	baseDir string
}

type LocalStorageConfig struct {
	BaseDir string
}

func NewLocalStorage(config LocalStorageConfig) (*LocalStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(config.BaseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base directory: %v", err)
	}

	return &LocalStorage{
		baseDir: config.BaseDir,
	}, nil
}

func (l *LocalStorage) UploadFile(ctx context.Context, file io.Reader, fileName string, contentType string) (string, error) {
	// Generate a unique filename
	uniqueName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(fileName))
	filePath := filepath.Join(l.baseDir, uniqueName)

	// Create the file
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer outFile.Close()

	// Copy the content
	if _, err := io.Copy(outFile, file); err != nil {
		// Clean up the file if copy fails
		os.Remove(filePath)
		return "", fmt.Errorf("failed to copy file content: %v", err)
	}

	return uniqueName, nil
}

func (l *LocalStorage) DownloadFile(ctx context.Context, fileName string) (io.ReadCloser, error) {
	filePath := filepath.Join(l.baseDir, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	return file, nil
}

func (l *LocalStorage) DeleteFile(ctx context.Context, fileName string) error {
	filePath := filepath.Join(l.baseDir, fileName)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}
	return nil
}

func (l *LocalStorage) GetFileURL(fileName string) string {
	return fmt.Sprintf("/api/v1/files/%s/download", fileName)
} 