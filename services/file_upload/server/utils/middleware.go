package utils

import (
	TokenGeneratorClient "github.com/DeebTibi/GoVault/services/token_generator/client"
	"github.com/gofiber/fiber/v2"
)

// Middleware to check for client token
func Authenticate(c *fiber.Ctx) error {
	clientToken := c.Get("Client-Token")
	userId := c.Get("User-ID")

	if userId == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing User-ID header")
	}

	if clientToken == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("Missing client token")
	}

	res, err := TokenGeneratorClient.NewTokenGeneratorClient().ValidateToken(userId, clientToken)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid client token")
	}

	if !res {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid client token")
	}

	// Proceed to the next middleware/handler
	return c.Next()
}
