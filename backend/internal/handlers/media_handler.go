package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/media-sync-player/backend/internal/models"
	"github.com/media-sync-player/backend/internal/services"
)

// MediaHandler handles HTTP requests related to media operations.
type MediaHandler struct {
	service *services.MediaService
}

// NewMediaHandler creates a new handler with the injected media service.
func NewMediaHandler(service *services.MediaService) *MediaHandler {
	return &MediaHandler{
		service: service,
	}
}

// sendJSONError is a helper function to send standard JSON errors.
func (h *MediaHandler) sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// ServeHTTP acts as a simple router for the /media prefix.
func (h *MediaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Strip the "/media" prefix to determine if there is an ID in the URL
	path := strings.TrimPrefix(r.URL.Path, "/media")

	if path == "" || path == "/" {
		// Root path requested: handle GET /media and POST /media
		switch r.Method {
		case http.MethodGet:
			h.getAllMedia(w, r)
		case http.MethodPost:
			h.createMedia(w, r)
		default:
			h.sendJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Path with ID requested (e.g., "/media/123")
	// Trim any remaining leading slashes to extract just the ID
	id := strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		h.getMediaByID(w, r, id)
	case http.MethodDelete:
		h.deleteMedia(w, r, id)
	default:
		h.sendJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// getAllMedia handles GET /media requests.
func (h *MediaHandler) getAllMedia(w http.ResponseWriter, r *http.Request) {
	mediaList, err := h.service.GetAllMedia()
	if err != nil {
		h.sendJSONError(w, "Failed to retrieve media", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mediaList)
}

// createMedia handles POST /media requests.
func (h *MediaHandler) createMedia(w http.ResponseWriter, r *http.Request) {
	var media models.Media
	if err := json.NewDecoder(r.Body).Decode(&media); err != nil {
		h.sendJSONError(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateMedia(media); err != nil {
		h.sendJSONError(w, "Failed to create media", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(media)
}

// getMediaByID handles GET /media/{id} requests.
func (h *MediaHandler) getMediaByID(w http.ResponseWriter, r *http.Request, id string) {
	media, err := h.service.GetMediaByID(id)
	if err != nil {
		// Map our specific repository error to a 404
		if err.Error() == "media not found" {
			h.sendJSONError(w, "media not found", http.StatusNotFound)
			return
		}
		h.sendJSONError(w, "Failed to retrieve media", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(media)
}

// deleteMedia handles DELETE /media/{id} requests.
func (h *MediaHandler) deleteMedia(w http.ResponseWriter, r *http.Request, id string) {
	err := h.service.DeleteMedia(id)
	if err != nil {
		// Map our specific repository error to a 404
		if err.Error() == "no media document found to delete" {
			h.sendJSONError(w, "media not found", http.StatusNotFound)
			return
		}
		h.sendJSONError(w, "Failed to delete media", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
