package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/media-sync-player/backend/internal/models"
	"github.com/media-sync-player/backend/internal/services"
)

// PlaylistHandler handles HTTP requests related to playlist operations.
type PlaylistHandler struct {
	service      *services.PlaylistService
	mediaService *services.MediaService
}

// NewPlaylistHandler creates a new handler with the injected playlist service.
func NewPlaylistHandler(service *services.PlaylistService, mediaService *services.MediaService) *PlaylistHandler {
	return &PlaylistHandler{
		service:      service,
		mediaService: mediaService,
	}
}

// sendJSONError is a helper function to send standard JSON errors.
func (h *PlaylistHandler) sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// ServeHTTP acts as a simple router for the /playlists prefix.
func (h *PlaylistHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Strip the "/playlists" prefix to determine if there is an ID in the URL
	path := strings.TrimPrefix(r.URL.Path, "/playlists")

	if path == "" || path == "/" {
		// Root path requested: handle GET /playlists and POST /playlists
		switch r.Method {
		case http.MethodGet:
			h.getAllPlaylists(w, r)
		case http.MethodPost:
			h.createPlaylist(w, r)
		default:
			h.sendJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Trim leading and trailing slashes
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	// If the path looks like /playlists/{id}/media
	if len(parts) == 2 && parts[1] == "media" {
		playlistID := parts[0]
		if r.Method == http.MethodPost {
			h.addMediaToPlaylist(w, r, playlistID)
			return
		}
		h.sendJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// If the path looks like /playlists/{id}
	if len(parts) == 1 {
		playlistID := parts[0]
		switch r.Method {
		case http.MethodGet:
			h.getPlaylistByID(w, r, playlistID)
		case http.MethodDelete:
			h.deletePlaylist(w, r, playlistID)
		default:
			h.sendJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// If it doesn't match known patterns
	h.sendJSONError(w, "Not Found", http.StatusNotFound)
}

// getAllPlaylists handles GET /playlists requests.
func (h *PlaylistHandler) getAllPlaylists(w http.ResponseWriter, r *http.Request) {
	playlists, err := h.service.GetAllPlaylists()
	if err != nil {
		h.sendJSONError(w, "Failed to retrieve playlists", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(playlists)
}

// createPlaylist handles POST /playlists requests.
func (h *PlaylistHandler) createPlaylist(w http.ResponseWriter, r *http.Request) {
	var playlist models.Playlist
	if err := json.NewDecoder(r.Body).Decode(&playlist); err != nil {
		h.sendJSONError(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.service.CreatePlaylist(playlist); err != nil {
		h.sendJSONError(w, "Failed to create playlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(playlist)
}

// getPlaylistByID handles GET /playlists/{id} requests.
func (h *PlaylistHandler) getPlaylistByID(w http.ResponseWriter, r *http.Request, id string) {
	playlist, err := h.service.GetPlaylistByID(id)
	if err != nil {
		if err.Error() == "playlist not found" {
			h.sendJSONError(w, "playlist not found", http.StatusNotFound)
			return
		}
		h.sendJSONError(w, "Failed to retrieve playlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(playlist)
}

// deletePlaylist handles DELETE /playlists/{id} requests.
func (h *PlaylistHandler) deletePlaylist(w http.ResponseWriter, r *http.Request, id string) {
	err := h.service.DeletePlaylist(id)
	if err != nil {
		if err.Error() == "no playlist document found to delete" || err.Error() == "playlist not found" {
			h.sendJSONError(w, "playlist not found", http.StatusNotFound)
			return
		}
		h.sendJSONError(w, "Failed to delete playlist", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

func (h *PlaylistHandler) addMediaToPlaylist(w http.ResponseWriter, r *http.Request, id string) {
	var req struct {
		MediaID string `json:"media_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendJSONError(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if req.MediaID == "" {
		h.sendJSONError(w, "media_id is required", http.StatusBadRequest)
		return
	}

	// Sanitize input to uppercase (media IDs are uppercase like 'M3', 'BLANK1')
	req.MediaID = strings.ToUpper(strings.TrimSpace(req.MediaID))

	// Validate media exists
	_, err := h.mediaService.GetMediaByID(req.MediaID)
	if err != nil {
		h.sendJSONError(w, "Media not found", http.StatusNotFound)
		return
	}

	// Add to playlist
	err = h.service.AddMediaToPlaylist(id, req.MediaID)
	if err != nil {
		if err.Error() == "playlist not found" {
			h.sendJSONError(w, "Playlist not found", http.StatusNotFound)
			return
		}
		h.sendJSONError(w, "Failed to add media to playlist", http.StatusInternalServerError)
		return
	}

	// Return updated playlist
	playlist, err := h.service.GetPlaylistByID(id)
	if err != nil {
		h.sendJSONError(w, "Failed to retrieve updated playlist", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(playlist)
}
