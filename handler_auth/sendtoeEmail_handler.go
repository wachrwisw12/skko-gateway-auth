package handler_auth

import (
	"skko-gateway-auth/db"
	"skko-gateway-auth/middleware"

	"github.com/gofiber/fiber/v2"
)

type RequestBody struct {
	Email string `json:"email"`
}

func SendToEmail(c *fiber.Ctx) error {
	var body RequestBody

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request json body",
		})
	}

	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM co_user WHERE email = ?)", body.Email).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "database error",
		})
	}

	if exists {
		otp, err := middleware.GenerateOTP(body.Email) // ✅ generate OTP 6 หลัก
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate OTP",
			})
		}

		err = middleware.SendOTPEmail(body.Email, otp)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to send email",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "ส่ง Email สำเร็จ",
		})
	} else {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "ไม่มี Email นี้ในระบบ",
		})
	}
}
