package handler

import (
	"github.com/ethandiaz/yad-back/internal/domain"
	"github.com/ethandiaz/yad-back/internal/repository"
	"github.com/ethandiaz/yad-back/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LocationHandler struct {
	service *service.LocationService
}

func NewLocationHandler(db *gorm.DB) *LocationHandler {
	repo := repository.NewLocationRepository(db)
	svc := service.NewLocationService(repo)
	return &LocationHandler{service: svc}
}

type CreateLocationRequest struct {
	Name             string  `json:"name" validate:"required"`
	Latitude         float64 `json:"latitude" validate:"required"`
	Longitude        float64 `json:"longitude" validate:"required"`
	FormattedAddress string  `json:"formattedAddress"`
}

type UpdateLocationRequest struct {
	Name             string `json:"name"`
	FormattedAddress string `json:"formattedAddress"`
}

func (h *LocationHandler) CreateLocation(c *fiber.Ctx) error {
	var req CreateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	location, err := h.service.CreateLocation(req.Name, req.Latitude, req.Longitude, req.FormattedAddress)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(location)
}

func (h *LocationHandler) GetLocation(c *fiber.Ctx) error {
	id := c.Params("id")
	location, err := h.service.GetLocation(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Location not found"})
	}

	return c.JSON(location)
}

func (h *LocationHandler) SearchLocations(c *fiber.Ctx) error {
	query := c.Query("query")
	limit := c.QueryInt("limit", 20)

	locations, err := h.service.SearchLocations(query, limit)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if locations == nil {
		locations = []*domain.Location{}
	}

	return c.JSON(fiber.Map{
		"total":     len(locations),
		"locations": locations,
	})
}

func (h *LocationHandler) FindNearbyLocations(c *fiber.Ctx) error {
	lat := c.QueryFloat("latitude")
	lon := c.QueryFloat("longitude")
	radius := c.QueryFloat("radius", 5)
	limit := c.QueryInt("limit", 20)

	if lat == 0 && lon == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Latitude and longitude are required"})
	}

	locations, err := h.service.FindNearbyLocations(lat, lon, radius, limit)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	if locations == nil {
		locations = []*domain.Location{}
	}

	return c.JSON(fiber.Map{
		"total":     len(locations),
		"locations": locations,
	})
}

func (h *LocationHandler) UpdateLocation(c *fiber.Ctx) error {
	id := c.Params("id")

	var req UpdateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.FormattedAddress != "" {
		updates["formattedAddress"] = req.FormattedAddress
	}

	location, err := h.service.UpdateLocation(id, updates)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(location)
}

func (h *LocationHandler) DeleteLocation(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.DeleteLocation(id); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Location not found"})
	}

	return c.SendStatus(204)
}
