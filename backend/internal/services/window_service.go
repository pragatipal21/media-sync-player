package services

import (
	"github.com/media-sync-player/backend/internal/models"
	"github.com/media-sync-player/backend/internal/repository"
)

// WindowService provides business logic for window operations.
type WindowService struct {
	repo *repository.WindowRepository
}

// NewWindowService creates a new window service with the injected repository.
func NewWindowService(repo *repository.WindowRepository) *WindowService {
	return &WindowService{
		repo: repo,
	}
}

// CreateWindow passes a window object to the repository for creation.
func (s *WindowService) CreateWindow(window models.Window) error {
	return s.repo.CreateWindow(window)
}

// GetWindowByID retrieves a window object from the repository by its ID.
func (s *WindowService) GetWindowByID(id string) (models.Window, error) {
	return s.repo.GetWindowByID(id)
}

// GetAllWindows retrieves all window objects from the repository.
func (s *WindowService) GetAllWindows() ([]models.Window, error) {
	return s.repo.GetAllWindows()
}

// DeleteWindow deletes a window object from the repository by its ID.
func (s *WindowService) DeleteWindow(id string) error {
	return s.repo.DeleteWindow(id)
}
