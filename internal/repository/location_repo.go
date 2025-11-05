package repository

import (
	"fmt"

	"github.com/ethandiaz/yad-back/internal/domain"
	"gorm.io/gorm"
)

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) domain.LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Create(location *domain.Location) error {
	return r.db.Create(location).Error
}

func (r *locationRepository) GetByID(id string) (*domain.Location, error) {
	var location domain.Location
	if err := r.db.First(&location, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *locationRepository) GetByName(name string, limit int) ([]*domain.Location, error) {
	var locations []*domain.Location
	if err := r.db.Where("name LIKE ?", "%"+name+"%").Limit(limit).Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

func (r *locationRepository) Update(location *domain.Location) error {
	return r.db.Save(location).Error
}

func (r *locationRepository) Delete(id string) error {
	return r.db.Delete(&domain.Location{}, "id = ?", id).Error
}

func (r *locationRepository) FindNearby(latitude, longitude, radiusKm float64, limit int) ([]*domain.Location, error) {
	var locations []*domain.Location

	// Haversine formula for distance calculation
	query := `
		SELECT * FROM locations
		WHERE ( 6371 * acos( cos( radians(?) ) * cos( radians( latitude ) ) 
			* cos( radians( longitude ) - radians(?) ) 
			+ sin( radians(?) ) * sin( radians( latitude ) ) ) ) <= ?
		ORDER BY 
			( 6371 * acos( cos( radians(?) ) * cos( radians( latitude ) ) 
			* cos( radians( longitude ) - radians(?) ) 
			+ sin( radians(?) ) * sin( radians( latitude ) ) ) ) ASC
		LIMIT ?
	`

	if err := r.db.Raw(query, latitude, longitude, latitude, radiusKm, latitude, longitude, latitude, limit).Scan(&locations).Error; err != nil {
		return nil, fmt.Errorf("failed to find nearby locations: %w", err)
	}
	return locations, nil
}
