package main

import (
	"context"
	"log"
	"os"
	"user-service/internal/handlers"
	"user-service/internal/repository"
	"user-service/internal/service"
	"user-service/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize MongoDB connection
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	// Get database
	db := mongoClient.Database("analyticsai")

	// Initialize Google Cloud Storage
	gcsConfig := storage.GCSConfig{
		ProjectID:       os.Getenv("GCS_PROJECT_ID"),
		BucketName:      os.Getenv("GCS_BUCKET_NAME"),
		CredentialsFile: os.Getenv("GCS_CREDENTIALS_FILE"),
	}

	gcsStorage, err := storage.NewGCSStorage(gcsConfig)
	if err != nil {
		log.Fatalf("Failed to initialize GCS storage: %v", err)
	}

	// Initialize repositories
	fileRepo := repository.NewFileRepository(db)

	// Initialize services
	fileService := service.NewFileService(fileRepo, gcsStorage)

	// Initialize handlers
	fileHandler := handlers.NewFileHandler(fileService)

	// Set up Gin router
	router := gin.Default()

	// Add middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Add authentication middleware
	router.Use(func(c *gin.Context) {
		// TODO: Implement proper authentication
		// For now, we'll just set a dummy user ID
		c.Set("user_id", uint(1))
		c.Next()
	})

	// API routes
	api := router.Group("/api/v1")
	{
		files := api.Group("/files")
		{
			files.POST("/upload", fileHandler.UploadFile)
			files.POST("/upload-url", fileHandler.UploadFileFromURL)
			files.GET("", fileHandler.ListFiles)
			files.DELETE("/:id", fileHandler.DeleteFile)
			files.PATCH("/:id/hide", fileHandler.HideFile)
			files.GET("/:id/download", fileHandler.DownloadFile)
		}
	}

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
