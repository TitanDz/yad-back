package repository

import (
	"github.com/ethandiaz/yad-back/internal/domain"
	"gorm.io/gorm"
)

type authTokenRepository struct {
	db *gorm.DB
}

func NewAuthTokenRepository(db *gorm.DB) domain.AuthTokenRepository {
	return &authTokenRepository{db: db}
}

func (r *authTokenRepository) Create(token *domain.AuthToken) error {
	return r.db.Create(token).Error
}

func (r *authTokenRepository) GetByToken(token string) (*domain.AuthToken, error) {
	var authToken domain.AuthToken
	if err := r.db.First(&authToken, "token = ?", token).Error; err != nil {
		return nil, err
	}
	return &authToken, nil
}

func (r *authTokenRepository) GetByUserID(userID string) ([]*domain.AuthToken, error) {
	var tokens []*domain.AuthToken
	if err := r.db.Find(&tokens, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *authTokenRepository) Delete(id string) error {
	return r.db.Delete(&domain.AuthToken{}, "id = ?", id).Error
}

func (r *authTokenRepository) DeleteByUserID(userID string) error {
	return r.db.Delete(&domain.AuthToken{}, "user_id = ?", userID).Error
}
