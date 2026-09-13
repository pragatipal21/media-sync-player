package services

import (
	"github.com/media-sync-player/backend/internal/models"
	"github.com/media-sync-player/backend/internal/repository"
)

// PlaylistService provides business logic for playlist operations.
type PlaylistService struct {
	repo *repository.PlaylistRepository
}

// NewPlaylistService creates a new playlist service with the injected repository.
func NewPlaylistService(repo *repository.PlaylistRepository) *PlaylistService {
	return &PlaylistService{
		repo: repo,
	}
}

// CreatePlaylist passes a playlist object to the repository for creation.
func (s *PlaylistService) CreatePlaylist(playlist models.Playlist) error {
	return s.repo.CreatePlaylist(playlist)
}

// GetPlaylistByID retrieves a playlist object from the repository by its ID.
func (s *PlaylistService) GetPlaylistByID(id string) (models.Playlist, error) {
	return s.repo.GetPlaylistByID(id)
}

// GetAllPlaylists retrieves all playlist objects from the repository.
func (s *PlaylistService) GetAllPlaylists() ([]models.Playlist, error) {
	return s.repo.GetAllPlaylists()
}

// DeletePlaylist deletes a playlist object from the repository by its ID.
func (s *PlaylistService) DeletePlaylist(id string) error {
	return s.repo.DeletePlaylist(id)
}

// AddMediaToPlaylist adds a media item to an existing playlist
func (s *PlaylistService) AddMediaToPlaylist(id string, mediaID string) error {
	return s.repo.AddMediaToPlaylist(id, mediaID)
}
