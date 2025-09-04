package handler_auth

import (
	"strings"

	"skko-gateway-auth/middleware"

	"github.com/gofiber/fiber/v2"
)

func RefreshTokenHandler(c *fiber.Ctx) error {
	var req struct {
		SessionID string `json:"session_id"`
	}
	if err := c.BodyParser(&req); err != nil || req.SessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	newAccess, _, err := middleware.RefreshSessionToken(req.SessionID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "refresh_failed"})
	}
	print(newAccess)
	print(req.SessionID)
	return c.JSON(fiber.Map{
		"accessToken": newAccess,
		"session_id":  req.SessionID,
	})
}

func VerifyTokenHandler(c *fiber.Ctx) error {
	// ดึง token จาก header
	auth := c.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
	}

	token := strings.TrimPrefix(auth, "Bearer ")

	isExpired := middleware.IsTokenExpired(token)

	return c.JSON(fiber.Map{"expired": isExpired})
}
