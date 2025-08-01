package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/wachrwisw12/corework-gateway-auth/handler"
)

func SetupAuth(auth fiber.Router) {
	auth.Post("/userlogin", handler.LoginHandler)
}
