package service

import (
	"fmt"

	"github.com/ethandiaz/yad-back/internal/domain"
)

type LocationService struct {
	repo domain.LocationRepository
}

func NewLocationService(repo domain.LocationRepository) *LocationService {
	return &LocationService{repo: repo}
}

func (s *LocationService) CreateLocation(name string, lat, lon float64, formattedAddress string) (*domain.Location, error) {
	if name == "" || lat < -90 || lat > 90 || lon < -180 || lon > 180 {
		return nil, fmt.Errorf("invalid location data")
	}

	location := domain.NewLocation(name, lat, lon, formattedAddress)
	if err := s.repo.Create(location); err != nil {
		return nil, fmt.Errorf("failed to create location: %w", err)
	}

	return location, nil
}

func (s *LocationService) GetLocation(id string) (*domain.Location, error) {
	location, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("location not found: %w", err)
	}
	return location, nil
}

func (s *LocationService) SearchLocations(query string, limit int) ([]*domain.Location, error) {
	if query == "" {
		return nil, fmt.Errorf("search query is required")
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	locations, err := s.repo.GetByName(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search locations: %w", err)
	}

	return locations, nil
}

func (s *LocationService) FindNearbyLocations(latitude, longitude, radiusKm float64, limit int) ([]*domain.Location, error) {
	if latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 {
		return nil, fmt.Errorf("invalid coordinates")
	}

	if radiusKm <= 0 {
		radiusKm = 5 // Default 5 km
	}

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	locations, err := s.repo.FindNearby(latitude, longitude, radiusKm, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearby locations: %w", err)
	}

	return locations, nil
}

func (s *LocationService) UpdateLocation(id string, updates map[string]interface{}) (*domain.Location, error) {
	location, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("location not found")
	}

	if name, ok := updates["name"].(string); ok && name != "" {
		location.Name = name
	}
	if formattedAddress, ok := updates["formattedAddress"].(string); ok {
		location.FormattedAddress = formattedAddress
	}

	if err := s.repo.Update(location); err != nil {
		return nil, fmt.Errorf("failed to update location: %w", err)
	}

	return location, nil
}

func (s *LocationService) DeleteLocation(id string) error {
	return s.repo.Delete(id)
}
