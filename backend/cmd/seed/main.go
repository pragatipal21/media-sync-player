package main

import (
	"log"
	"time"

	"github.com/media-sync-player/backend/internal/config"
	"github.com/media-sync-player/backend/internal/models"
	"github.com/media-sync-player/backend/internal/repository"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB using existing logic
	db, err := repository.ConnectMongoDB(cfg.MongoURI, cfg.MongoDatabase)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer db.Close()

	log.Println("MongoDB connected successfully for seeding.")

	// Instantiate repositories
	mediaRepo := repository.NewMediaRepository(db)
	playlistRepo := repository.NewPlaylistRepository(db)
	windowRepo := repository.NewWindowRepository(db)

	now := time.Now()

	// 1. Seed Media
	medias := []models.Media{
		{ID: "M1", Name: "Media 1", Type: "video", Duration: 60, FilePath: "media/m1.mp4", CreatedAt: now},
		{ID: "M2", Name: "Media 2", Type: "video", Duration: 90, FilePath: "media/m2.mp4", CreatedAt: now},
		{ID: "M3", Name: "Media 3", Type: "image", Duration: 30, FilePath: "media/m3.jpg", CreatedAt: now},
		{ID: "M4", Name: "Media 4", Type: "video", Duration: 120, FilePath: "media/m4.mp4", CreatedAt: now},
		{ID: "M5", Name: "Media 5", Type: "image", Duration: 45, FilePath: "media/m5.jpg", CreatedAt: now},
		{ID: "BLANK1", Name: "Blank Interval", Type: "blank", Duration: 10, FilePath: "", CreatedAt: now},
	}

	log.Println("Seeding Media...")
	for _, m := range medias {
		// Simple idempotency: attempt to delete if it exists, then create it fresh
		_ = mediaRepo.DeleteMedia(m.ID)
		if err := mediaRepo.CreateMedia(m); err != nil {
			log.Printf(" - Failed to seed media %s: %v\n", m.ID, err)
		} else {
			log.Printf(" - Seeded media: %s\n", m.ID)
		}
	}

	// 2. Seed Playlists
	playlists := []models.Playlist{
		{ID: "playlist-001", Name: "Window 1 Playlist", MediaIDs: []string{"M1", "BLANK1", "M2", "M3"}, CreatedAt: now},
		{ID: "playlist-002", Name: "Window 2 Playlist", MediaIDs: []string{"M2", "M4"}, CreatedAt: now},
		{ID: "playlist-003", Name: "Window 3 Playlist", MediaIDs: []string{"M1", "M5"}, CreatedAt: now},
	}

	log.Println("Seeding Playlists...")
	for _, p := range playlists {
		_ = playlistRepo.DeletePlaylist(p.ID)
		if err := playlistRepo.CreatePlaylist(p); err != nil {
			log.Printf(" - Failed to seed playlist %s: %v\n", p.ID, err)
		} else {
			log.Printf(" - Seeded playlist: %s\n", p.ID)
		}
	}

	// 3. Seed Windows
	windows := []models.Window{
		{ID: "window-001", Name: "Window 1", PlaylistID: "playlist-001", Status: "idle", CreatedAt: now},
		{ID: "window-002", Name: "Window 2", PlaylistID: "playlist-002", Status: "idle", CreatedAt: now},
		{ID: "window-003", Name: "Window 3", PlaylistID: "playlist-003", Status: "idle", CreatedAt: now},
	}

	log.Println("Seeding Windows...")
	for _, w := range windows {
		_ = windowRepo.DeleteWindow(w.ID)
		if err := windowRepo.CreateWindow(w); err != nil {
			log.Printf(" - Failed to seed window %s: %v\n", w.ID, err)
		} else {
			log.Printf(" - Seeded window: %s\n", w.ID)
		}
	}

	log.Println("Database seeding completed successfully.")
}
