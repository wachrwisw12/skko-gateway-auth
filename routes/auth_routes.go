package routes

import (
	"skko-gateway-auth/handler"

	"github.com/gofiber/fiber/v2"
)

func SetupAuth(auth fiber.Router) {
	auth.Post("/userlogin", handler.LoginHandler)
}
