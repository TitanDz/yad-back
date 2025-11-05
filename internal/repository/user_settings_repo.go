package repository

import (
	"github.com/ethandiaz/yad-back/internal/domain"
	"gorm.io/gorm"
)

type userSettingsRepository struct {
	db *gorm.DB
}

func NewUserSettingsRepository(db *gorm.DB) domain.UserSettingsRepository {
	return &userSettingsRepository{db: db}
}

func (r *userSettingsRepository) Create(settings *domain.UserSettings) error {
	return r.db.Create(settings).Error
}

func (r *userSettingsRepository) GetByUserID(userID string) (*domain.UserSettings, error) {
	var settings domain.UserSettings
	if err := r.db.First(&settings, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *userSettingsRepository) Update(settings *domain.UserSettings) error {
	return r.db.Save(settings).Error
}

func (r *userSettingsRepository) Delete(userID string) error {
	return r.db.Delete(&domain.UserSettings{}, "user_id = ?", userID).Error
}
