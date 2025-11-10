package service

import (
	"fmt"

	"github.com/ethandiaz/yad-back/internal/domain"
)

// NotificationService handles business logic for notifications
type NotificationService struct {
	repo domain.NotificationRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService(repo domain.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

// CreateNotification creates and saves a new notification
func (s *NotificationService) CreateNotification(userID, title, message, notificationType string, relatedMinyanID *string) (*domain.Notification, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}
	if title == "" || message == "" {
		return nil, fmt.Errorf("title and message are required")
	}

	notification := domain.NewNotification(userID, title, message, notificationType, relatedMinyanID)
	if err := s.repo.Create(notification); err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

// GetNotification retrieves a notification by ID
func (s *NotificationService) GetNotification(id string) (*domain.Notification, error) {
	notification, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	if notification == nil {
		return nil, fmt.Errorf("notification not found")
	}
	return notification, nil
}

// GetUserNotifications retrieves all notifications for a user
func (s *NotificationService) GetUserNotifications(userID string, limit, offset int) ([]*domain.Notification, error) {
	if userID == "" {
		return nil, fmt.Errorf("user_id is required")
	}

	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	notifications, err := s.repo.GetByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user notifications: %w", err)
	}

	if notifications == nil {
		notifications = []*domain.Notification{}
	}

	return notifications, nil
}

// GetUnreadCount returns the count of unread notifications
func (s *NotificationService) GetUnreadCount(userID string) (int64, error) {
	if userID == "" {
		return 0, fmt.Errorf("user_id is required")
	}

	count, err := s.repo.GetUnreadCount(userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(id string) error {
	if id == "" {
		return fmt.Errorf("notification_id is required")
	}

	// Verify notification exists
	notification, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	if notification == nil {
		return fmt.Errorf("notification not found")
	}

	if err := s.repo.MarkAsRead(id); err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	return nil
}

// MarkAllAsRead marks all notifications for a user as read
func (s *NotificationService) MarkAllAsRead(userID string) error {
	if userID == "" {
		return fmt.Errorf("user_id is required")
	}

	if err := s.repo.MarkAllAsRead(userID); err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	return nil
}

// DeleteNotification deletes a notification
func (s *NotificationService) DeleteNotification(id string) error {
	if id == "" {
		return fmt.Errorf("notification_id is required")
	}

	// Verify notification exists
	notification, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	if notification == nil {
		return fmt.Errorf("notification not found")
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	return nil
}

// SendNotification sends a notification to a user (business logic for notification trigger)
func (s *NotificationService) SendNotification(userID, title, message, notificationType string, relatedMinyanID *string) (*domain.Notification, error) {
	return s.CreateNotification(userID, title, message, notificationType, relatedMinyanID)
}
