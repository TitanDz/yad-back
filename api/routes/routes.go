package routes

import (
"github.com/ethandiaz/yad-back/internal/handler"
"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, h *handler.Handlers) {
	v1 := app.Group("/api/v1")

	// User routes
	users := v1.Group("/users")
	users.Post("/register", h.User.Register)
	users.Get("/:id", h.User.GetUser)
	users.Delete("/:id", h.User.DeleteUser)

	// Health check
	v1.Get("/health", func(c *fiber.Ctx) error {
return c.JSON(fiber.Map{"status": "healthy"})
})
}
