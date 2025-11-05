package service

import (
	"fmt"
	"time"

	"github.com/ethandiaz/yad-back/internal/domain"
)

type MinyanService struct {
	repo domain.MinyanRepository
}

func NewMinyanService(repo domain.MinyanRepository) *MinyanService {
	return &MinyanService{repo: repo}
}

func (s *MinyanService) CreateMinyan(userID, prayerType, date, timeStr, locationName string, lat, lon float64, notes string) (*domain.Minyan, error) {
	// Validate inputs
	if userID == "" || prayerType == "" || date == "" || timeStr == "" || locationName == "" {
		return nil, fmt.Errorf("missing required fields")
	}

	minyan := domain.NewMinyan(userID, prayerType, date, timeStr, locationName, lat, lon)
	minyan.Notes = notes

	if err := s.repo.Create(minyan); err != nil {
		return nil, fmt.Errorf("failed to create minyan: %w", err)
	}

	return minyan, nil
}

func (s *MinyanService) GetMinyan(id string) (*domain.Minyan, error) {
	minyan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("minyan not found: %w", err)
	}
	return minyan, nil
}

func (s *MinyanService) GetUserMinyans(userID string) ([]*domain.Minyan, error) {
	minyans, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user minyans: %w", err)
	}
	return minyans, nil
}

func (s *MinyanService) UpdateMinyan(id, userID string, updates map[string]interface{}) (*domain.Minyan, error) {
	minyan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("minyan not found")
	}

	// Check ownership
	if minyan.UserID != userID {
		return nil, fmt.Errorf("unauthorized: only creator can update minyan")
	}

	// Update allowed fields
	if prayerType, ok := updates["prayerType"].(string); ok && prayerType != "" {
		minyan.PrayerType = prayerType
	}
	if date, ok := updates["date"].(string); ok && date != "" {
		minyan.Date = date
	}
	if timeStr, ok := updates["time"].(string); ok && timeStr != "" {
		minyan.Time = timeStr
	}
	if locationName, ok := updates["locationName"].(string); ok && locationName != "" {
		minyan.LocationName = locationName
	}
	if lat, ok := updates["latitude"].(float64); ok {
		minyan.Latitude = lat
	}
	if lon, ok := updates["longitude"].(float64); ok {
		minyan.Longitude = lon
	}
	if notes, ok := updates["notes"].(string); ok {
		minyan.Notes = notes
	}

	minyan.UpdatedAt = time.Now()

	if err := s.repo.Update(minyan); err != nil {
		return nil, fmt.Errorf("failed to update minyan: %w", err)
	}

	return minyan, nil
}

func (s *MinyanService) DeleteMinyan(id, userID string) error {
	minyan, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("minyan not found")
	}

	// Check ownership
	if minyan.UserID != userID {
		return fmt.Errorf("unauthorized: only creator can delete minyan")
	}

	return s.repo.Delete(id)
}

func (s *MinyanService) PublishMinyan(id, userID string) (*domain.Minyan, error) {
	minyan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("minyan not found")
	}

	// Check ownership
	if minyan.UserID != userID {
		return nil, fmt.Errorf("unauthorized: only creator can publish minyan")
	}

	minyan.Status = "published"
	minyan.UpdatedAt = time.Now()

	if err := s.repo.Update(minyan); err != nil {
		return nil, fmt.Errorf("failed to publish minyan: %w", err)
	}

	return minyan, nil
}

func (s *MinyanService) SearchMinyans(prayerType, date string, limit, offset int) ([]*domain.Minyan, error) {
	// Validate limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	minyans, err := s.repo.Search(prayerType, date, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to search minyans: %w", err)
	}
	return minyans, nil
}

func (s *MinyanService) FindNearbyMinyans(latitude, longitude, radiusKm float64, limit int) ([]*domain.Minyan, error) {
	// Validate inputs
	if latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 {
		return nil, fmt.Errorf("invalid coordinates")
	}
	if radiusKm <= 0 {
		radiusKm = 5 // Default 5 km
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	minyans, err := s.repo.FindNearby(latitude, longitude, radiusKm, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to find nearby minyans: %w", err)
	}
	return minyans, nil
}

func (s *MinyanService) JoinMinyan(id string) (*domain.Minyan, error) {
	minyan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("minyan not found")
	}

	minyan.ParticipantCount++
	minyan.UpdatedAt = time.Now()

	if err := s.repo.Update(minyan); err != nil {
		return nil, fmt.Errorf("failed to join minyan: %w", err)
	}

	return minyan, nil
}

func (s *MinyanService) LeaveMinyan(id string) (*domain.Minyan, error) {
	minyan, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("minyan not found")
	}

	if minyan.ParticipantCount > 1 {
		minyan.ParticipantCount--
	}
	minyan.UpdatedAt = time.Now()

	if err := s.repo.Update(minyan); err != nil {
		return nil, fmt.Errorf("failed to leave minyan: %w", err)
	}

	return minyan, nil
}
