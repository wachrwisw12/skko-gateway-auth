package handler

import (
	"fmt"

	"skko-gateway-auth/middleware"
	"skko-gateway-auth/services"

	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(c *fiber.Ctx) error {
	var body LoginRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ข้อมูลไม่ถูกต้อง",
		})
	}

	user, err := services.LoginByUser(body.Username, body.Password)
	fmt.Println(user, err)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถเข้าสู่ระบบได้",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง",
		})
	}
	token, err := middleware.GenerateJWT(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถสร้าง token ได้",
		})
	}
	return c.JSON(fiber.Map{
		"message": "เข้าสู่ระบบสำเร็จ",
		"token":   token,
	})
}
