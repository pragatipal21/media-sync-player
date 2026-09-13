package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/media-sync-player/backend/internal/models"
)

// SyncHandler handles the simple in-memory synchronization logic.
type SyncHandler struct {
	mu    sync.Mutex
	state models.SyncState
}

// NewSyncHandler creates a new handler.
func NewSyncHandler() *SyncHandler {
	return &SyncHandler{
		state: models.SyncState{Active: false},
	}
}

// sendJSONError is a helper function to send standard JSON errors.
func (h *SyncHandler) sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// ServeHTTP acts as a simple router for the /sync prefix.
func (h *SyncHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/sync")

	if path != "" && path != "/" {
		h.sendJSONError(w, "Not Found", http.StatusNotFound)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getSyncState(w, r)
	case http.MethodPost:
		h.postSyncState(w, r)
	default:
		h.sendJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *SyncHandler) getSyncState(w http.ResponseWriter, r *http.Request) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Check if expired
	if h.state.Active && h.state.StartedAt != nil {
		expirationTime := h.state.StartedAt.Add(time.Duration(h.state.Duration) * time.Second)
		if time.Now().After(expirationTime) {
			h.state = models.SyncState{Active: false}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(h.state)
}

func (h *SyncHandler) postSyncState(w http.ResponseWriter, r *http.Request) {
	var req models.SyncRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendJSONError(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.MediaID == "" {
		h.sendJSONError(w, "media_id is required", http.StatusBadRequest)
		return
	}
	req.MediaID = strings.ToUpper(strings.TrimSpace(req.MediaID))

	if req.Duration <= 0 {
		h.sendJSONError(w, "duration must be greater than 0", http.StatusBadRequest)
		return
	}

	now := time.Now()

	h.mu.Lock()
	h.state = models.SyncState{
		Active:    true,
		MediaID:   req.MediaID,
		StartedAt: &now,
		Duration:  req.Duration,
	}
	newState := h.state
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newState)
}
