package repository

import (
	"github.com/ethandiaz/yad-back/internal/domain"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(db *gorm.DB) domain.NotificationRepository {
	return &NotificationRepository{db: db}
}

// Create saves a new notification to the database
func (r *NotificationRepository) Create(notification *domain.Notification) error {
	return r.db.Create(notification).Error
}

// GetByID retrieves a notification by ID
func (r *NotificationRepository) GetByID(id string) (*domain.Notification, error) {
	var notification domain.Notification
	err := r.db.Where("id = ?", id).First(&notification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &notification, nil
}

// GetByUserID retrieves all notifications for a user with pagination
func (r *NotificationRepository) GetByUserID(userID string, limit, offset int) ([]*domain.Notification, error) {
	var notifications []*domain.Notification
	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifications).
		Error
	if err != nil {
		return nil, err
	}
	return notifications, nil
}

// GetUnreadCount returns the count of unread notifications for a user
func (r *NotificationRepository) GetUnreadCount(userID string) (int64, error) {
	var count int64
	err := r.db.
		Model(&domain.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).
		Error
	return count, err
}

// MarkAsRead marks a single notification as read
func (r *NotificationRepository) MarkAsRead(id string) error {
	return r.db.Model(&domain.Notification{}).Where("id = ?", id).Update("is_read", true).Error
}

// MarkAllAsRead marks all notifications for a user as read
func (r *NotificationRepository) MarkAllAsRead(userID string) error {
	return r.db.
		Model(&domain.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).
		Error
}

// Delete removes a notification by ID
func (r *NotificationRepository) Delete(id string) error {
	return r.db.Delete(&domain.Notification{}, "id = ?", id).Error
}

// DeleteByUserID removes all notifications for a user
func (r *NotificationRepository) DeleteByUserID(userID string) error {
	return r.db.Delete(&domain.Notification{}, "user_id = ?", userID).Error
}

// Update updates an existing notification
func (r *NotificationRepository) Update(notification *domain.Notification) error {
	return r.db.Save(notification).Error
}
