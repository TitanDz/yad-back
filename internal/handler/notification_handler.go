package handler

import (
	"strconv"

	"github.com/ethandiaz/yad-back/internal/domain"
	"github.com/ethandiaz/yad-back/internal/repository"
	"github.com/ethandiaz/yad-back/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	service *service.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(db *gorm.DB) *NotificationHandler {
	repo := repository.NewNotificationRepository(db)
	svc := service.NewNotificationService(repo)
	return &NotificationHandler{service: svc}
}

// CreateNotificationRequest represents a request to create a notification
type CreateNotificationRequest struct {
	Title           string  `json:"title" validate:"required,min=1,max=255"`
	Message         string  `json:"message" validate:"required,min=1"`
	Type            string  `json:"type" validate:"required,oneof=minyan_alert minyan_update announcement"`
	RelatedMinyanID *string `json:"relatedMinyanId"`
}

// NotificationResponse represents a notification in the response
type NotificationResponse struct {
	ID              string  `json:"id"`
	UserID          string  `json:"userId"`
	Title           string  `json:"title"`
	Message         string  `json:"message"`
	Type            string  `json:"type"`
	RelatedMinyanID *string `json:"relatedMinyanId"`
	Icon            *string `json:"icon"`
	IsRead          bool    `json:"isRead"`
	CreatedAt       string  `json:"createdAt"`
}

// toNotificationResponse converts a notification domain model to a response
func toNotificationResponse(notification *domain.Notification) *NotificationResponse {
	return &NotificationResponse{
		ID:              notification.ID,
		UserID:          notification.UserID,
		Title:           notification.Title,
		Message:         notification.Message,
		Type:            notification.Type,
		RelatedMinyanID: notification.RelatedMinyanID,
		Icon:            notification.Icon,
		IsRead:          notification.IsRead,
		CreatedAt:       notification.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// CreateNotification handles POST /api/v1/notifications
func (h *NotificationHandler) CreateNotification(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	var req CreateNotificationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// TODO: Add request validation using validator package

	notification, err := h.service.CreateNotification(userID, req.Title, req.Message, req.Type, req.RelatedMinyanID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(toNotificationResponse(notification))
}

// GetNotifications handles GET /api/v1/notifications
func (h *NotificationHandler) GetNotifications(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	// Parse pagination parameters
	limit := 20
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset > 0 {
			offset = parsedOffset
		}
	}

	notifications, err := h.service.GetUserNotifications(userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	responses := make([]*NotificationResponse, len(notifications))
	for i, notification := range notifications {
		responses[i] = toNotificationResponse(notification)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":   responses,
		"count":  len(responses),
		"offset": offset,
		"limit":  limit,
	})
}

// GetUnreadCount handles GET /api/v1/notifications/unread-count
func (h *NotificationHandler) GetUnreadCount(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	count, err := h.service.GetUnreadCount(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"unreadCount": count,
	})
}

// MarkAsRead handles PUT /api/v1/notifications/:id/read
func (h *NotificationHandler) MarkAsRead(c *fiber.Ctx) error {
	notificationID := c.Params("id")

	if err := h.service.MarkAsRead(notificationID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

// MarkAllAsRead handles PUT /api/v1/notifications/read-all
func (h *NotificationHandler) MarkAllAsRead(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	if err := h.service.MarkAllAsRead(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

// DeleteNotification handles DELETE /api/v1/notifications/:id
func (h *NotificationHandler) DeleteNotification(c *fiber.Ctx) error {
	notificationID := c.Params("id")

	if err := h.service.DeleteNotification(notificationID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}
