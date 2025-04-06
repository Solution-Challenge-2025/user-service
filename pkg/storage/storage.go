package storage

import (
	"context"
	"io"
)

// Storage defines the interface for file storage operations
type Storage interface {
	// UploadFile uploads a file and returns a unique identifier for the file
	UploadFile(ctx context.Context, file io.Reader, fileName string, contentType string) (string, error)
	
	// DownloadFile retrieves a file by its identifier
	DownloadFile(ctx context.Context, fileName string) (io.ReadCloser, error)
	
	// DeleteFile removes a file by its identifier
	DeleteFile(ctx context.Context, fileName string) error
	
	// GetFileURL returns the URL for downloading a file
	GetFileURL(fileName string) string
} 