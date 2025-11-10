package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Place represents a saved place/venue by a user
type Place struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	UserID      string         `gorm:"index" json:"userId"`
	PlaceID     string         `json:"placeId"`
	Name        string         `json:"name"`
	Address     string         `json:"address"`
	Latitude    float64        `json:"latitude"`
	Longitude   float64        `json:"longitude"`
	PlaceType   *string        `json:"placeType"`
	PhoneNumber *string        `json:"phoneNumber"`
	Website     *string        `json:"website"`
	SavedAt     time.Time      `json:"savedAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// NewPlace creates a new place instance
func NewPlace(userID, placeID, name, address string, latitude, longitude float64) *Place {
	return &Place{
		ID:        uuid.New().String(),
		UserID:    userID,
		PlaceID:   placeID,
		Name:      name,
		Address:   address,
		Latitude:  latitude,
		Longitude: longitude,
		SavedAt:   time.Now(),
		UpdatedAt: time.Now(),
	}
}

// PlaceRepository defines data access operations for places
type PlaceRepository interface {
	Create(place *Place) error
	GetByID(id string) (*Place, error)
	GetByUserID(userID string, limit, offset int) ([]*Place, error)
	GetByPlaceID(userID, placeID string) (*Place, error)
	Delete(id string) error
	DeleteByUserID(userID string) error
	Update(place *Place) error
	IsSaved(userID, placeID string) (bool, error)
	Search(userID, query string, limit, offset int) ([]*Place, error)
}
