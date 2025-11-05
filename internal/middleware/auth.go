package middleware

import (
	"strings"

	"github.com/ethandiaz/yad-back/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	// Get Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"error": "Missing authorization header"})
	}

	// Extract token from "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid authorization header format"})
	}

	token := parts[1]

	// Verify token
	jwtManager := jwt.NewJWTManager()
	claims, err := jwtManager.VerifyToken(token)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Store claims in context for later use
	c.Locals("userID", claims.UserID)
	c.Locals("email", claims.Email)

	return c.Next()
}
