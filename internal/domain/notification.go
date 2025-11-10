package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification represents a user notification
type Notification struct {
	ID              string         `gorm:"primaryKey" json:"id"`
	UserID          string         `gorm:"index" json:"userId"`
	Title           string         `json:"title"`
	Message         string         `gorm:"type:text" json:"message"`
	Type            string         `json:"type"` // minyan_alert, minyan_update, announcement
	RelatedMinyanID *string        `json:"relatedMinyanId"`
	Icon            *string        `json:"icon"`
	IsRead          bool           `json:"isRead"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// NewNotification creates a new notification instance
func NewNotification(userID, title, message, notificationType string, relatedMinyanID *string) *Notification {
	return &Notification{
		ID:              uuid.New().String(),
		UserID:          userID,
		Title:           title,
		Message:         message,
		Type:            notificationType,
		RelatedMinyanID: relatedMinyanID,
		IsRead:          false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
}

// NotificationRepository defines data access operations for notifications
type NotificationRepository interface {
	Create(notification *Notification) error
	GetByID(id string) (*Notification, error)
	GetByUserID(userID string, limit, offset int) ([]*Notification, error)
	GetUnreadCount(userID string) (int64, error)
	MarkAsRead(id string) error
	MarkAllAsRead(userID string) error
	Delete(id string) error
	DeleteByUserID(userID string) error
	Update(notification *Notification) error
}
