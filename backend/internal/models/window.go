package models

import "time"

// Window represents a playback session or screen in the multi-window player.
type Window struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	PlaylistID string    `json:"playlist_id"`
	Status     string    `json:"status"` // e.g., "playing", "paused", "idle"
	CreatedAt  time.Time `json:"created_at"`
}
