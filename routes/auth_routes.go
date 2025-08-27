package routes

import (
	"skko-gateway-auth/handler_auth"
	handlerleave "skko-gateway-auth/handler_leave"
	"skko-gateway-auth/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupAuth(auth fiber.Router) {
	auth.Post("/loginByEmail", handler_auth.LoginHandler)
	auth.Post("/refresh", handler_auth.RefreshHandler)
	// auth.Post("/sendtoEmail", handler_auth.SendToEmail)
	auth.Post("/verify-otp", handler_auth.VerifyOtpHandler)
	auth.Post("/checkQrcode", handler_auth.CheckQrcodeHandler)
	// หน้าหลัก
	// Group routes ที่ต้อง login
	protected := auth.Group("/", middleware.JWTProtected())
	// ระบบ ลงเวลา
	protected.Post("/timekeeping", handler_auth.Timekeeping)

	// ระบบ ลา
	protected.Post("/getleave", handlerleave.GetleaveHandler)
}
