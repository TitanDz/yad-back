package handler

import (
	"strconv"

	"github.com/ethandiaz/yad-back/internal/domain"
	"github.com/ethandiaz/yad-back/internal/repository"
	"github.com/ethandiaz/yad-back/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PlaceHandler handles place/venue-related HTTP requests
type PlaceHandler struct {
	service *service.PlaceService
}

// NewPlaceHandler creates a new place handler
func NewPlaceHandler(db *gorm.DB) *PlaceHandler {
	repo := repository.NewPlaceRepository(db)
	svc := service.NewPlaceService(repo)
	return &PlaceHandler{service: svc}
}

// SavePlaceRequest represents a request to save a place
type SavePlaceRequest struct {
	PlaceID     string  `json:"placeId" validate:"required"`
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Address     string  `json:"address" validate:"required,min=1"`
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
	PlaceType   *string `json:"placeType"`
	PhoneNumber *string `json:"phoneNumber"`
	Website     *string `json:"website"`
}

// UpdatePlaceRequest represents a request to update a place
type UpdatePlaceRequest struct {
	Name        *string `json:"name"`
	Address     *string `json:"address"`
	PlaceType   *string `json:"placeType"`
	PhoneNumber *string `json:"phoneNumber"`
	Website     *string `json:"website"`
}

// PlaceResponse represents a place in the response
type PlaceResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"userId"`
	PlaceID     string  `json:"placeId"`
	Name        string  `json:"name"`
	Address     string  `json:"address"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	PlaceType   *string `json:"placeType"`
	PhoneNumber *string `json:"phoneNumber"`
	Website     *string `json:"website"`
	SavedAt     string  `json:"savedAt"`
}

// toPlaceResponse converts a place domain model to a response
func toPlaceResponse(place *domain.Place) *PlaceResponse {
	return &PlaceResponse{
		ID:          place.ID,
		UserID:      place.UserID,
		PlaceID:     place.PlaceID,
		Name:        place.Name,
		Address:     place.Address,
		Latitude:    place.Latitude,
		Longitude:   place.Longitude,
		PlaceType:   place.PlaceType,
		PhoneNumber: place.PhoneNumber,
		Website:     place.Website,
		SavedAt:     place.SavedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// SavePlace handles POST /api/v1/places/saved
func (h *PlaceHandler) SavePlace(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	var req SavePlaceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// TODO: Add request validation using validator package

	place, err := h.service.SavePlace(userID, req.PlaceID, req.Name, req.Address, req.Latitude, req.Longitude, req.PlaceType, req.PhoneNumber, req.Website)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(toPlaceResponse(place))
}

// GetSavedPlaces handles GET /api/v1/places/saved
func (h *PlaceHandler) GetSavedPlaces(c *fiber.Ctx) error {
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

	places, err := h.service.GetUserSavedPlaces(userID, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	responses := make([]*PlaceResponse, len(places))
	for i, place := range places {
		responses[i] = toPlaceResponse(place)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":   responses,
		"count":  len(responses),
		"offset": offset,
		"limit":  limit,
	})
}

// CheckSavedPlace handles GET /api/v1/places/saved/:placeId
func (h *PlaceHandler) CheckSavedPlace(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	placeID := c.Params("placeId")

	saved, err := h.service.IsPlaceSaved(userID, placeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"placeId": placeID,
		"isSaved": saved,
	})
}

// RemoveSavedPlace handles DELETE /api/v1/places/saved/:id
func (h *PlaceHandler) RemoveSavedPlace(c *fiber.Ctx) error {
	placeID := c.Params("id")

	if err := h.service.RemoveSavedPlace(placeID); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusOK)
}

// UpdatePlace handles PUT /api/v1/places/saved/:id
func (h *PlaceHandler) UpdatePlace(c *fiber.Ctx) error {
	placeID := c.Params("id")

	var req UpdatePlaceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	place, err := h.service.UpdatePlace(placeID, req.Name, req.Address, req.PlaceType, req.PhoneNumber, req.Website)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(toPlaceResponse(place))
}

// SearchPlaces handles GET /api/v1/places/search
func (h *PlaceHandler) SearchPlaces(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)
	query := c.Query("q", "")

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

	places, err := h.service.SearchPlaces(userID, query, limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	responses := make([]*PlaceResponse, len(places))
	for i, place := range places {
		responses[i] = toPlaceResponse(place)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":   responses,
		"count":  len(responses),
		"offset": offset,
		"limit":  limit,
		"query":  query,
	})
}
