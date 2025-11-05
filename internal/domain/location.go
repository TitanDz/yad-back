package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Location struct {
	ID               string         `gorm:"primaryKey" json:"id"`
	Name             string         `json:"name"`
	Latitude         float64        `json:"latitude"`
	Longitude        float64        `json:"longitude"`
	FormattedAddress string         `json:"formattedAddress"`
	PlaceID          string         `json:"placeId"` // Google Places ID
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func NewLocation(name string, lat, lon float64, formattedAddress string) *Location {
	return &Location{
		ID:               uuid.New().String(),
		Name:             name,
		Latitude:         lat,
		Longitude:        lon,
		FormattedAddress: formattedAddress,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

type LocationRepository interface {
	Create(location *Location) error
	GetByID(id string) (*Location, error)
	GetByName(name string, limit int) ([]*Location, error)
	Update(location *Location) error
	Delete(id string) error
	FindNearby(latitude, longitude, radiusKm float64, limit int) ([]*Location, error)
}
