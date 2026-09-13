package models

import "time"

// Playlist represents a collection of media items to be played.
type Playlist struct {
	ID        string    `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	MediaIDs  []string  `json:"media_ids" bson:"media_ids"` // Ordered list of media IDs
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
