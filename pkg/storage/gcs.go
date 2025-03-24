package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSConfig struct {
	ProjectID       string
	BucketName      string
	CredentialsFile string
}

type GCSStorage struct {
	client     *storage.Client
	bucketName string
}

func NewGCSStorage(config GCSConfig) (*GCSStorage, error) {
	ctx := context.Background()

	// Initialize GCS client
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(config.CredentialsFile))
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %v", err)
	}

	// Check if bucket exists
	bucket := client.Bucket(config.BucketName)
	_, err = bucket.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("bucket %s does not exist or is not accessible: %v", config.BucketName, err)
	}

	return &GCSStorage{
		client:     client,
		bucketName: config.BucketName,
	}, nil
}

func (g *GCSStorage) UploadFile(ctx context.Context, file io.Reader, fileName string, contentType string) (string, error) {
	// Generate a unique object name
	objectName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(fileName))

	// Get bucket and create object
	bucket := g.client.Bucket(g.bucketName)
	obj := bucket.Object(objectName)
	writer := obj.NewWriter(ctx)

	// Set content type
	writer.ContentType = contentType

	// Copy file content
	if _, err := io.Copy(writer, file); err != nil {
		return "", fmt.Errorf("failed to copy file content: %v", err)
	}

	// Close writer
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	return objectName, nil
}

func (g *GCSStorage) DownloadFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	bucket := g.client.Bucket(g.bucketName)
	obj := bucket.Object(objectName)
	reader, err := obj.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %v", err)
	}

	return reader, nil
}

func (g *GCSStorage) DeleteFile(ctx context.Context, objectName string) error {
	bucket := g.client.Bucket(g.bucketName)
	obj := bucket.Object(objectName)
	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete object: %v", err)
	}

	return nil
}

func (g *GCSStorage) GetFileURL(objectName string) string {
	return fmt.Sprintf("/api/v1/files/%s/download", objectName)
}
