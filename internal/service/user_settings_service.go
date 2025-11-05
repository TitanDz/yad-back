package service

import (
	"fmt"
	"time"

	"github.com/ethandiaz/yad-back/internal/domain"
)

type UserSettingsService struct {
	repo domain.UserSettingsRepository
}

func NewUserSettingsService(repo domain.UserSettingsRepository) *UserSettingsService {
	return &UserSettingsService{repo: repo}
}

func (s *UserSettingsService) CreateSettings(userID string) (*domain.UserSettings, error) {
	settings := domain.NewUserSettings(userID)
	if err := s.repo.Create(settings); err != nil {
		return nil, fmt.Errorf("failed to create settings: %w", err)
	}
	return settings, nil
}

func (s *UserSettingsService) GetSettings(userID string) (*domain.UserSettings, error) {
	settings, err := s.repo.GetByUserID(userID)
	if err != nil {
		// If settings don't exist, create defaults
		return s.CreateSettings(userID)
	}
	return settings, nil
}

func (s *UserSettingsService) UpdateSettings(userID string, updates map[string]interface{}) (*domain.UserSettings, error) {
	settings, err := s.repo.GetByUserID(userID)
	if err != nil {
		// If settings don't exist, create defaults first
		settings, err = s.CreateSettings(userID)
		if err != nil {
			return nil, err
		}
	}

	// Update allowed fields
	if notificationsEnabled, ok := updates["notificationsEnabled"].(bool); ok {
		settings.NotificationsEnabled = notificationsEnabled
	}
	if emailNotifications, ok := updates["emailNotifications"].(bool); ok {
		settings.EmailNotifications = emailNotifications
	}
	if visibilityMode, ok := updates["visibilityMode"].(string); ok && visibilityMode != "" {
		// Validate visibility mode
		if visibilityMode == "public" || visibilityMode == "private" || visibilityMode == "friends_only" {
			settings.VisibilityMode = visibilityMode
		}
	}
	if travelMode, ok := updates["travelMode"].(bool); ok {
		settings.TravelMode = travelMode
	}
	if language, ok := updates["language"].(string); ok && language != "" {
		settings.Language = language
	}
	if timeFormat, ok := updates["timeFormat"].(string); ok && timeFormat != "" {
		if timeFormat == "12h" || timeFormat == "24h" {
			settings.TimeFormat = timeFormat
		}
	}
	if theme, ok := updates["theme"].(string); ok && theme != "" {
		if theme == "light" || theme == "dark" {
			settings.Theme = theme
		}
	}

	settings.UpdatedAt = time.Now()

	if err := s.repo.Update(settings); err != nil {
		return nil, fmt.Errorf("failed to update settings: %w", err)
	}

	return settings, nil
}

func (s *UserSettingsService) UpdateSingleSetting(userID, key string, value interface{}) (*domain.UserSettings, error) {
	updates := map[string]interface{}{
		key: value,
	}
	return s.UpdateSettings(userID, updates)
}

func (s *UserSettingsService) DeleteSettings(userID string) error {
	return s.repo.Delete(userID)
}
