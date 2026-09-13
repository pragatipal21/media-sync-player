package services

import (
	"github.com/media-sync-player/backend/internal/models"
	"github.com/media-sync-player/backend/internal/repository"
)

// MediaService provides business logic for media operations.
type MediaService struct {
	repo *repository.MediaRepository
}

// NewMediaService creates a new media service with the injected repository.
func NewMediaService(repo *repository.MediaRepository) *MediaService {
	return &MediaService{
		repo: repo,
	}
}

// CreateMedia passes a media object to the repository for creation.
func (s *MediaService) CreateMedia(media models.Media) error {
	return s.repo.CreateMedia(media)
}

// GetMediaByID retrieves a media object from the repository by its ID.
func (s *MediaService) GetMediaByID(id string) (models.Media, error) {
	return s.repo.GetMediaByID(id)
}

// GetAllMedia retrieves all media objects from the repository.
func (s *MediaService) GetAllMedia() ([]models.Media, error) {
	return s.repo.GetAllMedia()
}

// DeleteMedia deletes a media object from the repository by its ID.
func (s *MediaService) DeleteMedia(id string) error {
	return s.repo.DeleteMedia(id)
}
