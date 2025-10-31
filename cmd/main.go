package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ethandiaz/yad-back/api/routes"
	"github.com/ethandiaz/yad-back/internal/config"
	"github.com/ethandiaz/yad-back/internal/database"
	"github.com/ethandiaz/yad-back/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/helmet/v2"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}
}

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "YAD API v1",
	})

	// Add security middleware
	app.Use(helmet.New())

	// Initialize handlers with database
	handlers := handler.NewHandlers(db)

	// Setup routes
	routes.SetupRoutes(app, handlers)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"time":   time.Now(),
		})
	})

	// Start server
	port := cfg.Port
	log.Printf("Starting YAD API server on port %s\n", port)
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
