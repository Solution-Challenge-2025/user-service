package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileStatus string

const (
	FileStatusActive    FileStatus = "active"
	FileStatusHidden    FileStatus = "hidden"
	FileStatusDeleted   FileStatus = "deleted"
	FileStatusAnalyzing FileStatus = "analyzing"
)

type File struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	UserID      uint               `bson:"user_id" json:"user_id"`
	Name        string             `bson:"name" json:"name"`
	OriginalURL string             `bson:"original_url,omitempty" json:"original_url,omitempty"`
	StorageKey  string             `bson:"storage_key" json:"storage_key"`
	Size        int64              `bson:"size" json:"size"`
	MimeType    string             `bson:"mime_type" json:"mime_type"`
	Status      FileStatus         `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type FileUploadRequest struct {
	Name        string `json:"name" binding:"required"`
	URL         string `json:"url,omitempty"`
	ContentType string `json:"content_type,omitempty"`
}

type FileResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Size        int64      `json:"size"`
	Status      FileStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	DownloadURL string     `json:"download_url,omitempty"`
}
