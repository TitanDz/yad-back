package handler

import (
	"github.com/ethandiaz/yad-back/internal/domain"
	"github.com/ethandiaz/yad-back/internal/repository"
	"github.com/ethandiaz/yad-back/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MinyanHandler struct {
	service *service.MinyanService
}

func NewMinyanHandler(db *gorm.DB) *MinyanHandler {
	repo := repository.NewMinyanRepository(db)
	svc := service.NewMinyanService(repo)
	return &MinyanHandler{service: svc}
}

type CreateMinyanRequest struct {
	PrayerType   string  `json:"prayerType" validate:"required"`
	Date         string  `json:"date" validate:"required"`
	Time         string  `json:"time" validate:"required"`
	LocationName string  `json:"locationName" validate:"required"`
	Latitude     float64 `json:"latitude" validate:"required"`
	Longitude    float64 `json:"longitude" validate:"required"`
	Notes        string  `json:"notes"`
}

type UpdateMinyanRequest struct {
	PrayerType   string   `json:"prayerType"`
	Date         string   `json:"date"`
	Time         string   `json:"time"`
	LocationName string   `json:"locationName"`
	Latitude     *float64 `json:"latitude"`
	Longitude    *float64 `json:"longitude"`
	Notes        string   `json:"notes"`
}

func (h *MinyanHandler) CreateMinyan(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	var req CreateMinyanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	minyan, err := h.service.CreateMinyan(
		userID, req.PrayerType, req.Date, req.Time,
		req.LocationName, req.Latitude, req.Longitude, req.Notes,
	)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(minyan)
}

func (h *MinyanHandler) GetMinyan(c *fiber.Ctx) error {
	id := c.Params("id")
	minyan, err := h.service.GetMinyan(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Minyan not found"})
	}

	return c.JSON(minyan)
}

func (h *MinyanHandler) GetUserMinyans(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)

	minyans, err := h.service.GetUserMinyans(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if minyans == nil {
		minyans = []*domain.Minyan{}
	}

	return c.JSON(fiber.Map{
		"total":   len(minyans),
		"minyans": minyans,
	})
}

func (h *MinyanHandler) UpdateMinyan(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("userID").(string)

	var req UpdateMinyanRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	updates := make(map[string]interface{})
	if req.PrayerType != "" {
		updates["prayerType"] = req.PrayerType
	}
	if req.Date != "" {
		updates["date"] = req.Date
	}
	if req.Time != "" {
		updates["time"] = req.Time
	}
	if req.LocationName != "" {
		updates["locationName"] = req.LocationName
	}
	if req.Latitude != nil {
		updates["latitude"] = *req.Latitude
	}
	if req.Longitude != nil {
		updates["longitude"] = *req.Longitude
	}
	if req.Notes != "" {
		updates["notes"] = req.Notes
	}

	minyan, err := h.service.UpdateMinyan(id, userID, updates)
	if err != nil {
		if err.Error() == "unauthorized: only creator can update minyan" {
			return c.Status(403).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(minyan)
}

func (h *MinyanHandler) DeleteMinyan(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("userID").(string)

	err := h.service.DeleteMinyan(id, userID)
	if err != nil {
		if err.Error() == "unauthorized: only creator can delete minyan" {
			return c.Status(403).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}

func (h *MinyanHandler) PublishMinyan(c *fiber.Ctx) error {
	id := c.Params("id")
	userID := c.Locals("userID").(string)

	minyan, err := h.service.PublishMinyan(id, userID)
	if err != nil {
		if err.Error() == "unauthorized: only creator can publish minyan" {
			return c.Status(403).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(minyan)
}

func (h *MinyanHandler) SearchMinyans(c *fiber.Ctx) error {
	prayerType := c.Query("prayerType")
	date := c.Query("date")
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	minyans, err := h.service.SearchMinyans(prayerType, date, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	if minyans == nil {
		minyans = []*domain.Minyan{}
	}

	return c.JSON(fiber.Map{
		"total":   len(minyans),
		"minyans": minyans,
	})
}

func (h *MinyanHandler) FindNearbyMinyans(c *fiber.Ctx) error {
	lat := c.QueryFloat("latitude")
	lon := c.QueryFloat("longitude")
	radius := c.QueryFloat("radius", 5)
	limit := c.QueryInt("limit", 20)

	if lat == 0 && lon == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Latitude and longitude are required"})
	}

	minyans, err := h.service.FindNearbyMinyans(lat, lon, radius, limit)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if minyans == nil {
		minyans = []*domain.Minyan{}
	}

	return c.JSON(fiber.Map{
		"total":   len(minyans),
		"minyans": minyans,
	})
}

func (h *MinyanHandler) JoinMinyan(c *fiber.Ctx) error {
	id := c.Params("id")

	minyan, err := h.service.JoinMinyan(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(minyan)
}

func (h *MinyanHandler) LeaveMinyan(c *fiber.Ctx) error {
	id := c.Params("id")

	minyan, err := h.service.LeaveMinyan(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(minyan)
}
