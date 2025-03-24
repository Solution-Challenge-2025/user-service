package repository

import (
	"context"
	"log"
	"time"
	"user-service/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FileRepository struct {
	collection *mongo.Collection
}

func NewFileRepository(db *mongo.Database) *FileRepository {
	return &FileRepository{
		collection: db.Collection("files"),
	}
}

func (r *FileRepository) Create(ctx context.Context, file *models.File) error {
	log.Printf("[FileRepository.Create] Starting file creation")

	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()

	// Create a new ObjectID
	file.ID = primitive.NewObjectID()
	log.Printf("[FileRepository.Create] Generated new ID: %s", file.ID.Hex())

	result, err := r.collection.InsertOne(ctx, file)
	if err != nil {
		log.Printf("[FileRepository.Create] Failed to insert file: %v", err)
		return err
	}
	log.Printf("Results: %s", result)

	log.Printf("[FileRepository.Create] File created successfully with ID: %s", file.ID.Hex())
	return nil
}

func (r *FileRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.File, error) {
	log.Printf("[FileRepository.GetByID] Fetching file with ID: %s", id.Hex())

	var file models.File
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&file)
	if err != nil {
		log.Printf("[FileRepository.GetByID] Failed to fetch file: %v", err)
		return nil, err
	}
	log.Printf("[FileRepository.GetByID] Successfully fetched file: %s", file.ID.Hex())
	return &file, nil
}

func (r *FileRepository) GetByUserID(ctx context.Context, userID uint) ([]models.File, error) {
	log.Printf("[FileRepository.GetByUserID] Fetching files for user: %d", userID)

	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		log.Printf("[FileRepository.GetByUserID] Failed to fetch files: %v", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var files []models.File
	if err := cursor.All(ctx, &files); err != nil {
		log.Printf("[FileRepository.GetByUserID] Failed to decode files: %v", err)
		return nil, err
	}
	log.Printf("[FileRepository.GetByUserID] Successfully fetched %d files", len(files))
	return files, nil
}

func (r *FileRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.FileStatus) error {
	log.Printf("[FileRepository.UpdateStatus] Updating status for file: %s to: %s", id.Hex(), status)

	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{
			"$set": bson.M{
				"status":     status,
				"updated_at": time.Now(),
			},
		},
	)
	if err != nil {
		log.Printf("[FileRepository.UpdateStatus] Failed to update status: %v", err)
		return err
	}
	log.Printf("[FileRepository.UpdateStatus] Successfully updated status")
	return nil
}

func (r *FileRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	log.Printf("[FileRepository.Delete] Deleting file: %s", id.Hex())

	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Printf("[FileRepository.Delete] Failed to delete file: %v", err)
		return err
	}
	log.Printf("[FileRepository.Delete] Successfully deleted file")
	return nil
}
