package routes

import (
	"github.com/ethandiaz/yad-back/internal/handler"
	"github.com/ethandiaz/yad-back/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, h *handler.Handlers) {
	v1 := app.Group("/api/v1")

	// User routes
	users := v1.Group("/users")
	users.Post("/register", h.User.Register)
	users.Get("/:id", h.User.GetUser)
	users.Delete("/:id", h.User.DeleteUser)

	// Auth routes
	auth := v1.Group("/auth")
	auth.Post("/login", h.Auth.Login)
	auth.Post("/refresh", h.Auth.RefreshToken)
	auth.Post("/logout", h.Auth.Logout)
	auth.Post("/verify", h.Auth.VerifyToken)

	// Minyan routes (protected)
	minyans := v1.Group("/minyans")
	minyans.Use(middleware.AuthMiddleware)
	minyans.Post("", h.Minyan.CreateMinyan)
	minyans.Get("/my", h.Minyan.GetUserMinyans)
	minyans.Get("/search", h.Minyan.SearchMinyans)
	minyans.Get("/nearby", h.Minyan.FindNearbyMinyans)
	minyans.Get("/:id", h.Minyan.GetMinyan)
	minyans.Put("/:id", h.Minyan.UpdateMinyan)
	minyans.Delete("/:id", h.Minyan.DeleteMinyan)
	minyans.Post("/:id/publish", h.Minyan.PublishMinyan)
	minyans.Post("/:id/join", h.Minyan.JoinMinyan)
	minyans.Post("/:id/leave", h.Minyan.LeaveMinyan)

	// Location routes (protected)
	locations := v1.Group("/locations")
	locations.Use(middleware.AuthMiddleware)
	locations.Post("", h.Location.CreateLocation)
	locations.Get("/search", h.Location.SearchLocations)
	locations.Get("/nearby", h.Location.FindNearbyLocations)
	locations.Get("/:id", h.Location.GetLocation)
	locations.Put("/:id", h.Location.UpdateLocation)
	locations.Delete("/:id", h.Location.DeleteLocation)

	// User settings routes (protected)
	settings := v1.Group("/users/:userId/settings")
	settings.Use(middleware.AuthMiddleware)
	settings.Get("", h.Settings.GetSettings)
	settings.Put("", h.Settings.UpdateSettings)
	settings.Put("/:key", h.Settings.UpdateSingleSetting)

	// Notification routes (protected)
	notifications := v1.Group("/notifications")
	notifications.Use(middleware.AuthMiddleware)
	notifications.Post("", h.Notification.CreateNotification)
	notifications.Get("", h.Notification.GetNotifications)
	notifications.Get("/unread-count", h.Notification.GetUnreadCount)
	notifications.Put("/:id/read", h.Notification.MarkAsRead)
	notifications.Put("/read-all", h.Notification.MarkAllAsRead)
	notifications.Delete("/:id", h.Notification.DeleteNotification)

	// Place routes (protected)
	places := v1.Group("/places")
	places.Use(middleware.AuthMiddleware)
	places.Post("/saved", h.Place.SavePlace)
	places.Get("/saved", h.Place.GetSavedPlaces)
	places.Get("/saved/:placeId", h.Place.CheckSavedPlace)
	places.Delete("/saved/:id", h.Place.RemoveSavedPlace)
	places.Put("/saved/:id", h.Place.UpdatePlace)
	places.Get("/search", h.Place.SearchPlaces)

	// Health check
	v1.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "healthy"})
	})
}
