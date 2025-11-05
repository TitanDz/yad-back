package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSettings struct {
	ID                   string         `gorm:"primaryKey" json:"id"`
	UserID               string         `gorm:"uniqueIndex" json:"userId"`
	NotificationsEnabled bool           `json:"notificationsEnabled"`
	EmailNotifications   bool           `json:"emailNotifications"`
	VisibilityMode       string         `json:"visibilityMode"` // public, private, friends_only
	TravelMode           bool           `json:"travelMode"`
	Language             string         `json:"language"`   // en, es, he, etc.
	TimeFormat           string         `json:"timeFormat"` // 12h, 24h
	Theme                string         `json:"theme"`      // light, dark
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`
}

func NewUserSettings(userID string) *UserSettings {
	return &UserSettings{
		ID:                   uuid.New().String(),
		UserID:               userID,
		NotificationsEnabled: true,
		EmailNotifications:   true,
		VisibilityMode:       "public",
		TravelMode:           false,
		Language:             "en",
		TimeFormat:           "12h",
		Theme:                "light",
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

type UserSettingsRepository interface {
	Create(settings *UserSettings) error
	GetByUserID(userID string) (*UserSettings, error)
	Update(settings *UserSettings) error
	Delete(userID string) error
}
