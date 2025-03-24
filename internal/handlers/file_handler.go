package handlers

import (
	"log"
	"net/http"
	"user-service/internal/models"
	"user-service/internal/service"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fileService *service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get file from request"})
		return
	}

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
		return
	}
	defer src.Close()

	// Upload file
	fileRecord, err := h.fileService.UploadFile(c.Request.Context(), userID.(uint), src, file.Filename, file.Header.Get("Content-Type"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, fileRecord)
}

func (h *FileHandler) UploadFileFromURL(c *gin.Context) {
	log.Printf("[UploadFileFromURL] Starting request processing")

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("[UploadFileFromURL] User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	log.Printf("[UploadFileFromURL] User ID found: %v", userID)

	var req models.FileUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("[UploadFileFromURL] Failed to bind JSON request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	log.Printf("[UploadFileFromURL] Request body parsed successfully - URL: %s, Name: %s", req.URL, req.Name)

	fileRecord, err := h.fileService.UploadFileFromURL(c.Request.Context(), userID.(uint), req.URL, req.Name)
	if err != nil {
		log.Printf("[UploadFileFromURL] Service error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("[UploadFileFromURL] File uploaded successfully - ID: %d, Name: %s", fileRecord.ID, fileRecord.Name)

	c.JSON(http.StatusCreated, fileRecord)
}

func (h *FileHandler) ListFiles(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	files, err := h.fileService.ListUserFiles(c.Request.Context(), userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, files)
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	log.Printf("[DeleteFile] Starting file deletion")

	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("[DeleteFile] User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	log.Printf("[DeleteFile] User ID found: %v", userID)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("[DeleteFile] Invalid file ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	if err := h.fileService.DeleteFile(c.Request.Context(), id); err != nil {
		log.Printf("[DeleteFile] Failed to delete file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[DeleteFile] Successfully deleted file")
	c.Status(http.StatusNoContent)
}

func (h *FileHandler) HideFile(c *gin.Context) {
	log.Printf("[HideFile] Starting file hide operation")

	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("[HideFile] User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	log.Printf("[HideFile] User ID found: %v", userID)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("[HideFile] Invalid file ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	if err := h.fileService.HideFile(c.Request.Context(), id); err != nil {
		log.Printf("[HideFile] Failed to hide file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[HideFile] Successfully hid file")
	c.Status(http.StatusNoContent)
}

func (h *FileHandler) DownloadFile(c *gin.Context) {
	log.Printf("[DownloadFile] Starting file download")

	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("[DownloadFile] User ID not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	log.Printf("[DownloadFile] User ID found: %v", userID)

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		log.Printf("[DownloadFile] Invalid file ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file ID"})
		return
	}

	file, err := h.fileService.GetFile(c.Request.Context(), id)
	if err != nil {
		log.Printf("[DownloadFile] Failed to fetch file: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	reader, err := h.fileService.DownloadFile(c.Request.Context(), id)
	if err != nil {
		log.Printf("[DownloadFile] Failed to download file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer reader.Close()

	log.Printf("[DownloadFile] Successfully downloaded file")
	c.Header("Content-Disposition", "attachment; filename="+file.Name)
	c.Header("Content-Type", file.MimeType)
	c.DataFromReader(http.StatusOK, -1, file.MimeType, reader, nil)
}
