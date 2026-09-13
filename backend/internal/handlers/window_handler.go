package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/media-sync-player/backend/internal/models"
	"github.com/media-sync-player/backend/internal/services"
)

// WindowHandler handles HTTP requests related to window operations.
type WindowHandler struct {
	service *services.WindowService
}

// NewWindowHandler creates a new handler with the injected window service.
func NewWindowHandler(service *services.WindowService) *WindowHandler {
	return &WindowHandler{
		service: service,
	}
}

// sendJSONError is a helper function to send standard JSON errors.
func (h *WindowHandler) sendJSONError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// ServeHTTP acts as a simple router for the /windows prefix.
func (h *WindowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Strip the "/windows" prefix to determine if there is an ID in the URL
	path := strings.TrimPrefix(r.URL.Path, "/windows")

	if path == "" || path == "/" {
		// Root path requested: handle GET /windows and POST /windows
		switch r.Method {
		case http.MethodGet:
			h.getAllWindows(w, r)
		case http.MethodPost:
			h.createWindow(w, r)
		default:
			h.sendJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// Path with ID requested (e.g., "/windows/123")
	// Trim any remaining leading slashes to extract just the ID
	id := strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodGet:
		h.getWindowByID(w, r, id)
	case http.MethodDelete:
		h.deleteWindow(w, r, id)
	default:
		h.sendJSONError(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// getAllWindows handles GET /windows requests.
func (h *WindowHandler) getAllWindows(w http.ResponseWriter, r *http.Request) {
	windows, err := h.service.GetAllWindows()
	if err != nil {
		h.sendJSONError(w, "Failed to retrieve windows", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(windows)
}

// createWindow handles POST /windows requests.
func (h *WindowHandler) createWindow(w http.ResponseWriter, r *http.Request) {
	var window models.Window
	if err := json.NewDecoder(r.Body).Decode(&window); err != nil {
		h.sendJSONError(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateWindow(window); err != nil {
		h.sendJSONError(w, "Failed to create window", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(window)
}

// getWindowByID handles GET /windows/{id} requests.
func (h *WindowHandler) getWindowByID(w http.ResponseWriter, r *http.Request, id string) {
	window, err := h.service.GetWindowByID(id)
	if err != nil {
		if err.Error() == "window not found" {
			h.sendJSONError(w, "window not found", http.StatusNotFound)
			return
		}
		h.sendJSONError(w, "Failed to retrieve window", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(window)
}

// deleteWindow handles DELETE /windows/{id} requests.
func (h *WindowHandler) deleteWindow(w http.ResponseWriter, r *http.Request, id string) {
	err := h.service.DeleteWindow(id)
	if err != nil {
		if err.Error() == "no window document found to delete" || err.Error() == "window not found" {
			h.sendJSONError(w, "window not found", http.StatusNotFound)
			return
		}
		h.sendJSONError(w, "Failed to delete window", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
