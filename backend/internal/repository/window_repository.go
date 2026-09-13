package repository

import (
	"context"
	"errors"
	"time"

	"github.com/media-sync-player/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WindowRepository struct {
	collection *mongo.Collection
}

// NewWindowRepository creates a new repository to interact with the windows collection.
func NewWindowRepository(db *MongoDB) *WindowRepository {
	return &WindowRepository{
		collection: db.Database.Collection("windows"),
	}
}

// CreateWindow inserts a new window document into the collection.
func (r *WindowRepository) CreateWindow(window models.Window) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, window)
	return err
}

// GetWindowByID retrieves a window document by its string ID.
func (r *WindowRepository) GetWindowByID(id string) (models.Window, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var window models.Window
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&window)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return window, errors.New("window not found")
		}
		return window, err
	}

	return window, nil
}

// GetAllWindows retrieves all window documents from the collection.
func (r *WindowRepository) GetAllWindows() ([]models.Window, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	// Ensure the cursor is closed when we're done
	defer cursor.Close(ctx)

	var windows []models.Window
	// Decode all documents directly into the slice
	if err = cursor.All(ctx, &windows); err != nil {
		return nil, err
	}

	// If no documents were found, return an empty array instead of nil
	if windows == nil {
		windows = []models.Window{}
	}

	return windows, nil
}

// DeleteWindow deletes a window document by its string ID.
func (r *WindowRepository) DeleteWindow(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no window document found to delete")
	}

	return nil
}
