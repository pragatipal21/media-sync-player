package models

import "time"

// SyncRequest is the payload expected when creating a new sync state.
type SyncRequest struct {
	MediaID  string `json:"media_id"`
	Duration int    `json:"duration"` // in seconds
}

// SyncState represents the current global synchronization state.
type SyncState struct {
	Active    bool       `json:"active"`
	MediaID   string     `json:"media_id,omitempty"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	Duration  int        `json:"duration,omitempty"` // in seconds
}
