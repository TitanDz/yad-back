package repository

import (
	"github.com/ethandiaz/yad-back/internal/domain"
	"gorm.io/gorm"
)

type PlaceRepository struct {
	db *gorm.DB
}

// NewPlaceRepository creates a new place repository
func NewPlaceRepository(db *gorm.DB) domain.PlaceRepository {
	return &PlaceRepository{db: db}
}

// Create saves a new place to the database
func (r *PlaceRepository) Create(place *domain.Place) error {
	return r.db.Create(place).Error
}

// GetByID retrieves a place by ID
func (r *PlaceRepository) GetByID(id string) (*domain.Place, error) {
	var place domain.Place
	err := r.db.Where("id = ?", id).First(&place).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &place, nil
}

// GetByUserID retrieves all saved places for a user with pagination
func (r *PlaceRepository) GetByUserID(userID string, limit, offset int) ([]*domain.Place, error) {
	var places []*domain.Place
	err := r.db.
		Where("user_id = ?", userID).
		Order("saved_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&places).
		Error
	if err != nil {
		return nil, err
	}
	return places, nil
}

// GetByPlaceID retrieves a saved place by user ID and place ID
func (r *PlaceRepository) GetByPlaceID(userID, placeID string) (*domain.Place, error) {
	var place domain.Place
	err := r.db.
		Where("user_id = ? AND place_id = ?", userID, placeID).
		First(&place).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &place, nil
}

// Delete removes a place by ID
func (r *PlaceRepository) Delete(id string) error {
	return r.db.Delete(&domain.Place{}, "id = ?", id).Error
}

// DeleteByUserID removes all saved places for a user
func (r *PlaceRepository) DeleteByUserID(userID string) error {
	return r.db.Delete(&domain.Place{}, "user_id = ?", userID).Error
}

// Update updates an existing place
func (r *PlaceRepository) Update(place *domain.Place) error {
	return r.db.Save(place).Error
}

// IsSaved checks if a place is saved by a user
func (r *PlaceRepository) IsSaved(userID, placeID string) (bool, error) {
	var count int64
	err := r.db.
		Model(&domain.Place{}).
		Where("user_id = ? AND place_id = ?", userID, placeID).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Search searches for saved places by name or other criteria
func (r *PlaceRepository) Search(userID, query string, limit, offset int) ([]*domain.Place, error) {
	var places []*domain.Place
	err := r.db.
		Where("user_id = ? AND (name LIKE ? OR address LIKE ?)", userID, "%"+query+"%", "%"+query+"%").
		Order("saved_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&places).
		Error
	if err != nil {
		return nil, err
	}
	return places, nil
}
