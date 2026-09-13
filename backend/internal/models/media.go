package models

import "time"

// Media represents a single media item (e.g., a video or audio file).
type Media struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	FilePath  string    `json:"file_path"`
	Type      string    `json:"type"`
	Duration  int       `json:"duration"` // Duration in seconds
	CreatedAt time.Time `json:"created_at"`
}
