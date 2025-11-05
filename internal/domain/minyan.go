package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Minyan struct {
	ID               string         `gorm:"primaryKey" json:"id"`
	UserID           string         `gorm:"index" json:"userId"`
	PrayerType       string         `json:"prayerType"` // Shacharit, Mincha, Maariv
	Date             string         `json:"date"`       // YYYY-MM-DD
	Time             string         `json:"time"`       // HH:MM
	LocationName     string         `json:"locationName"`
	Latitude         float64        `json:"latitude"`
	Longitude        float64        `json:"longitude"`
	Notes            string         `json:"notes"`
	Status           string         `json:"status"` // draft, published, cancelled
	ParticipantCount int            `json:"participantCount"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func NewMinyan(userID, prayerType, date, timeStr, locationName string, lat, lon float64) *Minyan {
	return &Minyan{
		ID:               uuid.New().String(),
		UserID:           userID,
		PrayerType:       prayerType,
		Date:             date,
		Time:             timeStr,
		LocationName:     locationName,
		Latitude:         lat,
		Longitude:        lon,
		Status:           "draft",
		ParticipantCount: 1,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

type MinyanRepository interface {
	Create(minyan *Minyan) error
	GetByID(id string) (*Minyan, error)
	GetByUserID(userID string) ([]*Minyan, error)
	Update(minyan *Minyan) error
	Delete(id string) error
	Search(prayerType, date string, limit, offset int) ([]*Minyan, error)
	FindNearby(latitude, longitude, radiusKm float64, limit int) ([]*Minyan, error)
}
