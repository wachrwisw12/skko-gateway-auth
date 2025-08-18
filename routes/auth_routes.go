package routes

import (
	"skko-gateway-auth/handler_auth"

	"github.com/gofiber/fiber/v2"
)

func SetupAuth(auth fiber.Router) {
	auth.Post("/loginByEmail", handler_auth.LoginHandler)
	auth.Post("/sendtoEmail", handler_auth.SendToEmail)
	auth.Post("/verifyOtp", handler_auth.VerifyOtpHandler)
}
