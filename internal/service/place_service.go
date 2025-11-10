package service

import (
	"fmt"

	"github.com/ethandiaz/yad-back/internal/domain"
)

// PlaceService handles business logic for places/saved venues
type PlaceService struct {
	repo domain.PlaceRepository
}

// NewPlaceService creates a new place service
func NewPlaceService(repo domain.PlaceRepository) *PlaceService {
	return &PlaceService{repo: repo}
}

// SavePlace saves a new place for a user
func (s *PlaceService) SavePlace(userID, placeID, name, address string, latitude, longitude float64, placeType, phoneNumber, website *string) (*domain.Place, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if name == "" || address == "" {
		return nil, fmt.Errorf("name and address are required")
	}

	// Check if place is already saved
	existing, err := s.repo.GetByPlaceID(userID, placeID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing place: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("place is already saved")
	}

	place := domain.NewPlace(userID, placeID, name, address, latitude, longitude)
	place.PlaceType = placeType
	place.PhoneNumber = phoneNumber
	place.Website = website

	if err := s.repo.Create(place); err != nil {
		return nil, fmt.Errorf("failed to save place: %w", err)
	}

	return place, nil
}

// GetSavedPlace retrieves a saved place by ID
func (s *PlaceService) GetSavedPlace(id string) (*domain.Place, error) {
	place, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get place: %w", err)
	}
	if place == nil {
		return nil, fmt.Errorf("place not found")
	}
	return place, nil
}

// GetUserSavedPlaces retrieves all saved places for a user with pagination
func (s *PlaceService) GetUserSavedPlaces(userID string, limit, offset int) ([]*domain.Place, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	places, err := s.repo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user saved places: %w", err)
	}

	if places == nil {
		places = []*domain.Place{}
	}

	return places, nil
}

// IsPlaceSaved checks if a place is saved by a user
func (s *PlaceService) IsPlaceSaved(userID, placeID string) (bool, error) {
	if userID == "" || placeID == "" {
		return false, fmt.Errorf("user_id and place_id are required")
	}

	saved, err := s.repo.IsSaved(userID, placeID)
	if err != nil {
		return false, fmt.Errorf("failed to check if place is saved: %w", err)
	}

	return saved, nil
}

// RemoveSavedPlace removes a saved place
func (s *PlaceService) RemoveSavedPlace(id string) error {
	if id == "" {
		return fmt.Errorf("place_id is required")
	}

	// Verify place exists
	place, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to remove place: %w", err)
	}
	if place == nil {
		return fmt.Errorf("place not found")
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to remove place: %w", err)
	}

	return nil
}

// UpdatePlace updates place information
func (s *PlaceService) UpdatePlace(id string, name, address *string, placeType, phoneNumber, website *string) (*domain.Place, error) {
	if id == "" {
		return nil, fmt.Errorf("place_id is required")
	}

	place, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to update place: %w", err)
	}
	if place == nil {
		return nil, fmt.Errorf("place not found")
	}

	// Update fields if provided
	if name != nil && *name != "" {
		place.Name = *name
	}
	if address != nil && *address != "" {
		place.Address = *address
	}
	if placeType != nil {
		place.PlaceType = placeType
	}
	if phoneNumber != nil {
		place.PhoneNumber = phoneNumber
	}
	if website != nil {
		place.Website = website
	}

	if err := s.repo.Update(place); err != nil {
		return nil, fmt.Errorf("failed to update place: %w", err)
	}

	return place, nil
}

// SearchPlaces searches for saved places by query
func (s *PlaceService) SearchPlaces(userID, query string, limit, offset int) ([]*domain.Place, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	places, err := s.repo.Search(userID, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search places: %w", err)
	}

	if places == nil {
		places = []*domain.Place{}
	}

	return places, nil
}
