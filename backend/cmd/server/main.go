package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/media-sync-player/backend/internal/config"
	"github.com/media-sync-player/backend/internal/handlers"
	"github.com/media-sync-player/backend/internal/repository"
	"github.com/media-sync-player/backend/internal/services"
)

// corsMiddleware is a simple standard-library wrapper to allow cross-origin requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	db, err := repository.ConnectMongoDB(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Close()

	log.Println("MongoDB connected successfully.")

	// Wire up Media dependencies (MongoDB -> Repository -> Service -> Handler)
	mediaRepo := repository.NewMediaRepository(db)
	mediaService := services.NewMediaService(mediaRepo)
	mediaHandler := handlers.NewMediaHandler(mediaService)

	// Wire up Playlist dependencies
	playlistRepo := repository.NewPlaylistRepository(db)
	playlistService := services.NewPlaylistService(playlistRepo)
	playlistHandler := handlers.NewPlaylistHandler(playlistService, mediaService)

	// Wire up Window dependencies
	windowRepo := repository.NewWindowRepository(db)
	windowService := services.NewWindowService(windowRepo)
	windowHandler := handlers.NewWindowHandler(windowService)

	// Wire up Sync dependencies (in-memory, no service/repo needed)
	syncHandler := handlers.NewSyncHandler()

	// Use http.DefaultServeMux implicitly by calling http.Handle
	// Register Media routes
	http.Handle("/media", mediaHandler)
	http.Handle("/media/", mediaHandler)

	// Register Playlist routes
	http.Handle("/playlists", playlistHandler)
	http.Handle("/playlists/", playlistHandler)

	// Register Window routes
	http.Handle("/windows", windowHandler)
	http.Handle("/windows/", windowHandler)

	// Register Sync routes
	http.Handle("/sync", syncHandler)
	http.Handle("/sync/", syncHandler)

	// Health endpoint
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]string{"status": "ok"}
		json.NewEncoder(w).Encode(response)
	})

	port := ":8080"
	fmt.Printf("Server starting on port %s...\n", port)

	// Wrap the default mux with our CORS middleware
	handler := corsMiddleware(http.DefaultServeMux)

	if err := http.ListenAndServe(port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
