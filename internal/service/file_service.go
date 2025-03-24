package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/pkg/storage"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileService struct {
	repo    *repository.FileRepository
	storage *storage.GCSStorage
}

func NewFileService(repo *repository.FileRepository, storage *storage.GCSStorage) *FileService {
	return &FileService{
		repo:    repo,
		storage: storage,
	}
}

func (s *FileService) UploadFile(ctx context.Context, userID uint, file io.Reader, fileName string, contentType string) (*models.File, error) {
	log.Printf("[UploadFile] Starting file upload - UserID: %d, FileName: %s, ContentType: %s", userID, fileName, contentType)

	// Upload file to storage
	storageKey, err := s.storage.UploadFile(ctx, file, fileName, contentType)
	if err != nil {
		log.Printf("[UploadFile] Failed to upload file to storage: %v", err)
		return nil, fmt.Errorf("failed to upload file to storage: %v", err)
	}
	log.Printf("[UploadFile] File uploaded to storage successfully - StorageKey: %s", storageKey)

	// Create file record in database
	fileRecord := &models.File{
		UserID:     userID,
		Name:       fileName,
		StorageKey: storageKey,
		MimeType:   contentType,
		Status:     models.FileStatusActive,
	}

	if err := s.repo.Create(ctx, fileRecord); err != nil {
		log.Printf("[UploadFile] Failed to create file record in database: %v", err)
		// Cleanup storage if database operation fails
		_ = s.storage.DeleteFile(ctx, storageKey)
		return nil, fmt.Errorf("failed to create file record: %v", err)
	}
	log.Printf("[UploadFile] File record created successfully - ID: %d", fileRecord.ID)

	return fileRecord, nil
}

func (s *FileService) UploadFileFromURL(ctx context.Context, userID uint, url string, fileName string) (*models.File, error) {
	log.Printf("[UploadFileFromURL] Starting URL file upload - UserID: %d, URL: %s, FileName: %s", userID, url, fileName)

	// Download file from URL
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[UploadFileFromURL] Failed to download file from URL: %v", err)
		return nil, fmt.Errorf("failed to download file from URL: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[UploadFileFromURL] Failed to download file: HTTP %d", resp.StatusCode)
		return nil, fmt.Errorf("failed to download file: HTTP %d", resp.StatusCode)
	}
	log.Printf("[UploadFileFromURL] File downloaded successfully from URL - ContentType: %s", resp.Header.Get("Content-Type"))

	// Upload file to storage
	return s.UploadFile(ctx, userID, resp.Body, fileName, resp.Header.Get("Content-Type"))
}

func (s *FileService) GetFile(ctx context.Context, id primitive.ObjectID) (*models.File, error) {
	log.Printf("[FileService.GetFile] Fetching file with ID: %s", id.Hex())
	return s.repo.GetByID(ctx, id)
}

func (s *FileService) ListUserFiles(ctx context.Context, userID uint) ([]models.File, error) {
	log.Printf("[FileService.ListUserFiles] Fetching files for user: %d", userID)
	return s.repo.GetByUserID(ctx, userID)
}

func (s *FileService) DeleteFile(ctx context.Context, id primitive.ObjectID) error {
	log.Printf("[FileService.DeleteFile] Deleting file: %s", id.Hex())

	file, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("[FileService.DeleteFile] Failed to fetch file: %v", err)
		return err
	}

	// Soft delete in database
	if err := s.repo.UpdateStatus(ctx, id, models.FileStatusDeleted); err != nil {
		log.Printf("[FileService.DeleteFile] Failed to update file status: %v", err)
		return fmt.Errorf("failed to update file status: %v", err)
	}

	// Delete from storage
	if err := s.storage.DeleteFile(ctx, file.StorageKey); err != nil {
		log.Printf("[FileService.DeleteFile] Failed to delete file from storage: %v", err)
		return fmt.Errorf("failed to delete file from storage: %v", err)
	}

	log.Printf("[FileService.DeleteFile] Successfully deleted file")
	return nil
}

func (s *FileService) HideFile(ctx context.Context, id primitive.ObjectID) error {
	log.Printf("[FileService.HideFile] Hiding file: %s", id.Hex())
	return s.repo.UpdateStatus(ctx, id, models.FileStatusHidden)
}

func (s *FileService) DownloadFile(ctx context.Context, id primitive.ObjectID) (io.ReadCloser, error) {
	log.Printf("[FileService.DownloadFile] Downloading file: %s", id.Hex())

	file, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("[FileService.DownloadFile] Failed to fetch file: %v", err)
		return nil, err
	}

	return s.storage.DownloadFile(ctx, file.StorageKey)
}
