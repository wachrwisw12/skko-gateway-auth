package routes

import (
	"skko-gateway-auth/compservice"
	"skko-gateway-auth/handler_auth"
	"skko-gateway-auth/middleware"
	"skko-gateway-auth/services"
	"skko-gateway-auth/timestamp"

	"github.com/gofiber/fiber/v2"
)

func SetupAuth(auth fiber.Router) {
	auth.Post("/verify-otp", handler_auth.VerifyOtpHandler)
	auth.Post("/checkQrcode", handler_auth.CheckQrcodeHandler)
	auth.Post("/refresh-token", handler_auth.RefreshTokenHandler)
	auth.Post("/verify-token", handler_auth.VerifyTokenHandler)

	// comp-service
	auth.Post("/send-system", compservice.AgentSystem)
	auth.Post("/heartbeat", compservice.AgentHeartbeat)

	// หน้าหลัก
	// Group routes ที่ต้อง login
	protected := auth.Group("/", middleware.JWTProtected())
	// หน้าหลัก
	protected.Get("/home", services.HomeApp)
	// ระบบ ลงเวลา
	protected.Post("/hometimeStamp", timestamp.TimestampHome)
	protected.Post("/checkin-timeStamp", timestamp.TimestampCheckIn)
	protected.Get("/timekeeping-history", timestamp.TimeStampHistory)
	// protected.Get("/testtime", timestamp.TimeTest)

	// ระบบ ลา
	// protected.Post("/getleave", handlerleave.GetleaveHandler)
}
