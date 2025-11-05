package handler

import (
	"github.com/ethandiaz/yad-back/internal/repository"
	"github.com/ethandiaz/yad-back/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserSettingsHandler struct {
	service *service.UserSettingsService
}

func NewUserSettingsHandler(db *gorm.DB) *UserSettingsHandler {
	repo := repository.NewUserSettingsRepository(db)
	svc := service.NewUserSettingsService(repo)
	return &UserSettingsHandler{service: svc}
}

type UpdateSettingsRequest struct {
	NotificationsEnabled *bool  `json:"notificationsEnabled"`
	EmailNotifications   *bool  `json:"emailNotifications"`
	VisibilityMode       string `json:"visibilityMode"`
	TravelMode           *bool  `json:"travelMode"`
	Language             string `json:"language"`
	TimeFormat           string `json:"timeFormat"`
	Theme                string `json:"theme"`
}

func (h *UserSettingsHandler) GetSettings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	settings, err := h.service.GetSettings(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(settings)
}

func (h *UserSettingsHandler) UpdateSettings(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req UpdateSettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	updates := make(map[string]interface{})
	if req.NotificationsEnabled != nil {
		updates["notificationsEnabled"] = *req.NotificationsEnabled
	}
	if req.EmailNotifications != nil {
		updates["emailNotifications"] = *req.EmailNotifications
	}
	if req.VisibilityMode != "" {
		updates["visibilityMode"] = req.VisibilityMode
	}
	if req.TravelMode != nil {
		updates["travelMode"] = *req.TravelMode
	}
	if req.Language != "" {
		updates["language"] = req.Language
	}
	if req.TimeFormat != "" {
		updates["timeFormat"] = req.TimeFormat
	}
	if req.Theme != "" {
		updates["theme"] = req.Theme
	}

	settings, err := h.service.UpdateSettings(userID, updates)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(settings)
}

func (h *UserSettingsHandler) UpdateSingleSetting(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	key := c.Params("key")

	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	value, ok := req["value"]
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Missing value field"})
	}

	settings, err := h.service.UpdateSingleSetting(userID, key, value)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(settings)
}
