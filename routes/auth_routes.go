package routes

import (
	"skko-gateway-auth/handler_auth"
	handlerleave "skko-gateway-auth/handler_leave"
	"skko-gateway-auth/middleware"
	"skko-gateway-auth/timestamp"

	"github.com/gofiber/fiber/v2"
)

func SetupAuth(auth fiber.Router) {
	auth.Post("/verify-otp", handler_auth.VerifyOtpHandler)
	auth.Post("/checkQrcode", handler_auth.CheckQrcodeHandler)
	auth.Post("/refresh-token", handler_auth.RefreshTokenHandler)
	// หน้าหลัก
	// Group routes ที่ต้อง login
	protected := auth.Group("/", middleware.JWTProtected())
	// ระบบ ลงเวลา
	protected.Post("/hometimeStamp", timestamp.TimestampHome)

	// ระบบ ลา
	protected.Post("/getleave", handlerleave.GetleaveHandler)
}
