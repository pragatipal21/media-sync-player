package repository

import (
	"context"
	"errors"
	"time"

	"github.com/media-sync-player/backend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlaylistRepository struct {
	collection *mongo.Collection
}

// NewPlaylistRepository creates a new repository to interact with the playlists collection.
func NewPlaylistRepository(db *MongoDB) *PlaylistRepository {
	return &PlaylistRepository{
		collection: db.Database.Collection("playlists"),
	}
}

// CreatePlaylist inserts a new playlist document into the collection.
func (r *PlaylistRepository) CreatePlaylist(playlist models.Playlist) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, playlist)
	return err
}

// GetPlaylistByID retrieves a playlist document by its string ID.
func (r *PlaylistRepository) GetPlaylistByID(id string) (models.Playlist, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var playlist models.Playlist
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&playlist)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return playlist, errors.New("playlist not found")
		}
		return playlist, err
	}

	return playlist, nil
}

// GetAllPlaylists retrieves all playlist documents from the collection.
func (r *PlaylistRepository) GetAllPlaylists() ([]models.Playlist, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	// Ensure the cursor is closed when we're done
	defer cursor.Close(ctx)

	var playlists []models.Playlist
	// Decode all documents directly into the slice
	if err = cursor.All(ctx, &playlists); err != nil {
		return nil, err
	}

	// If no documents were found, return an empty array instead of nil
	if playlists == nil {
		playlists = []models.Playlist{}
	}

	return playlists, nil
}

// DeletePlaylist deletes a playlist document by its string ID.
func (r *PlaylistRepository) DeletePlaylist(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no playlist document found to delete")
	}

	return nil
}

// AddMediaToPlaylist adds a media ID to the end of a playlist's media_ids array.
func (r *PlaylistRepository) AddMediaToPlaylist(id string, mediaID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"id": id},
		bson.M{"$push": bson.M{"media_ids": mediaID}},
	)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("playlist not found")
	}

	return nil
}
