package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthToken struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	UserID    string         `gorm:"index" json:"userId"`
	Token     string         `gorm:"unique;index" json:"token"`
	Type      string         `json:"type"` // "access" or "refresh"
	ExpiresAt time.Time      `json:"expiresAt"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func NewAuthToken(userID, token, tokenType string, expiresAt time.Time) *AuthToken {
	return &AuthToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		Token:     token,
		Type:      tokenType,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

type AuthTokenRepository interface {
	Create(token *AuthToken) error
	GetByToken(token string) (*AuthToken, error)
	GetByUserID(userID string) ([]*AuthToken, error)
	Delete(id string) error
	DeleteByUserID(userID string) error
}
