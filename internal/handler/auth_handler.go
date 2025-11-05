package handler

import (
	"github.com/ethandiaz/yad-back/internal/repository"
	"github.com/ethandiaz/yad-back/internal/service"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewAuthTokenRepository(db)
	svc := service.NewAuthService(userRepo, tokenRepo)
	return &AuthHandler{service: svc}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}

type LogoutRequest struct {
	UserID string `json:"userId" validate:"required"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// TODO: Add request validation
	authResponse, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(authResponse)
}

func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	var req RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	authResponse, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(authResponse)
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var req LogoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := h.service.Logout(req.UserID); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to logout"})
	}

	return c.SendStatus(200)
}

func (h *AuthHandler) VerifyToken(c *fiber.Ctx) error {
	// Get token from Authorization header
	token := c.Get("Authorization")
	if token == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Missing authorization token"})
	}

	// Remove "Bearer " prefix if present
	if len(token) > 7 && token[:7] == "Bearer " {
		token = token[7:]
	}

	claims, err := h.service.VerifyToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}

	return c.Status(200).JSON(fiber.Map{
		"valid":  true,
		"userId": claims.UserID,
		"email":  claims.Email,
	})
}
