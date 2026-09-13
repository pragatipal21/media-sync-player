package repository

import (
	"context"
	"errors"
	"time"

	"github.com/media-sync-player/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MediaRepository struct {
	collection *mongo.Collection
}

// NewMediaRepository creates a new repository to interact with the media collection.
func NewMediaRepository(db *MongoDB) *MediaRepository {
	return &MediaRepository{
		collection: db.Database.Collection("media"),
	}
}

// CreateMedia inserts a new media document into the collection.
func (r *MediaRepository) CreateMedia(media models.Media) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, media)
	return err
}

// GetMediaByID retrieves a media document by its string ID.
func (r *MediaRepository) GetMediaByID(id string) (models.Media, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var media models.Media
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&media)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return media, errors.New("media not found")
		}
		return media, err
	}

	return media, nil
}

// GetAllMedia retrieves all media documents from the collection.
func (r *MediaRepository) GetAllMedia() ([]models.Media, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	// Ensure the cursor is closed when we're done
	defer cursor.Close(ctx)

	var mediaList []models.Media
	// Decode all documents directly into the slice
	if err = cursor.All(ctx, &mediaList); err != nil {
		return nil, err
	}

	// If no documents were found, return an empty array instead of nil
	if mediaList == nil {
		mediaList = []models.Media{}
	}

	return mediaList, nil
}

// DeleteMedia deletes a media document by its string ID.
func (r *MediaRepository) DeleteMedia(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no media document found to delete")
	}

	return nil
}
