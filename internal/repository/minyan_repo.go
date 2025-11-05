package repository

import (
	"fmt"

	"github.com/ethandiaz/yad-back/internal/domain"
	"gorm.io/gorm"
)

type minyanRepository struct {
	db *gorm.DB
}

func NewMinyanRepository(db *gorm.DB) domain.MinyanRepository {
	return &minyanRepository{db: db}
}

func (r *minyanRepository) Create(minyan *domain.Minyan) error {
	return r.db.Create(minyan).Error
}

func (r *minyanRepository) GetByID(id string) (*domain.Minyan, error) {
	var minyan domain.Minyan
	if err := r.db.First(&minyan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &minyan, nil
}

func (r *minyanRepository) GetByUserID(userID string) ([]*domain.Minyan, error) {
	var minyans []*domain.Minyan
	if err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&minyans).Error; err != nil {
		return nil, err
	}
	return minyans, nil
}

func (r *minyanRepository) Update(minyan *domain.Minyan) error {
	return r.db.Save(minyan).Error
}

func (r *minyanRepository) Delete(id string) error {
	return r.db.Delete(&domain.Minyan{}, "id = ?", id).Error
}

func (r *minyanRepository) Search(prayerType, date string, limit, offset int) ([]*domain.Minyan, error) {
	var minyans []*domain.Minyan
	query := r.db.Where("status = ?", "published")

	if prayerType != "" {
		query = query.Where("prayer_type = ?", prayerType)
	}
	if date != "" {
		query = query.Where("date = ?", date)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&minyans).Error; err != nil {
		return nil, err
	}
	return minyans, nil
}

func (r *minyanRepository) FindNearby(latitude, longitude, radiusKm float64, limit int) ([]*domain.Minyan, error) {
	var minyans []*domain.Minyan

	// Haversine formula for distance calculation
	// Using SQL to calculate distance and filter by radius
	query := `
		SELECT * FROM minyans 
		WHERE status = 'published'
		AND ( 6371 * acos( cos( radians(?) ) * cos( radians( latitude ) ) 
			* cos( radians( longitude ) - radians(?) ) 
			+ sin( radians(?) ) * sin( radians( latitude ) ) ) ) <= ?
		ORDER BY created_at DESC
		LIMIT ?
	`

	if err := r.db.Raw(query, latitude, longitude, latitude, radiusKm, limit).Scan(&minyans).Error; err != nil {
		return nil, fmt.Errorf("failed to find nearby minyans: %w", err)
	}
	return minyans, nil
}
