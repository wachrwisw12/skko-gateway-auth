package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	api := app.Group("/api", handler) // /api

	v1 := api.Group("/v1")   // /api/v1
	v1.Get("/list", handler) // /api/v1/list
	// v1.Get("/user", handler)        // /api/v1/user

	// v2 := api.Group("/v2", handler) // /api/v2
	// v2.Get("/list", handler)        // /api/v2/list
	// v2.Get("/user", handler)        // /api/v2/user

	log.Fatal(app.Listen(":3000"))
}

func handler(c *fiber.Ctx) error {
	return c.SendString("tesกtห")
}
