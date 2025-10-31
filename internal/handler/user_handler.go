package handler

import (
"github.com/ethandiaz/yad-back/internal/repository"
"github.com/ethandiaz/yad-back/internal/service"
"github.com/gofiber/fiber/v2"
"gorm.io/gorm"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	repo := repository.NewUserRepository(db)
	svc := service.NewUserService(repo)
	return &UserHandler{service: svc}
}

type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user, err := h.service.Register(req.Email, req.FirstName, req.LastName, req.Password)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(UserResponse{
ID:        user.ID,
Email:     user.Email,
FirstName: user.FirstName,
LastName:  user.LastName,
})
}

func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.service.GetUser(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(UserResponse{
ID:        user.ID,
Email:     user.Email,
FirstName: user.FirstName,
LastName:  user.LastName,
})
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteUser(id); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.SendStatus(204)
}
