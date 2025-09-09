package services

import "github.com/gofiber/fiber/v2"

func HomeApp(c *fiber.Ctx) error {
	return c.JSON("home ok")
}
