package handler_auth

import (
	"skko-gateway-auth/middleware"
	"skko-gateway-auth/models"
	"skko-gateway-auth/services"

	"github.com/gofiber/fiber/v2"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func LoginHandler(c *fiber.Ctx) error {
	var body LoginRequest

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ข้อมูลไม่ถูกต้อง",
		})
	}

	user, err := services.LoginByEmail(body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "ไม่สามารถเข้าสู่ระบบได้",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "ชื่อผู้ใช้หรือรหัสผ่านไม่ถูกต้อง",
		})
	}
	accessToken, err := middleware.GenerateAccessToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ไม่สามารถสร้าง token ได้",
		})
	}
	refreshToken, err := middleware.GenerateRefreshToken(*user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "ไม่สามารถสร้าง refresh token ได้"})
	}
	return c.JSON(fiber.Map{
		"message":      "เข้าสู่ระบบสำเร็จแล้ว",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"userinfo":     user,
	})
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func RefreshHandler(c *fiber.Ctx) error {
	var body RefreshRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	claims, err := middleware.ParseJWT(body.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid refresh token"})
	}

	// สร้าง access token ใหม่จาก claims
	user := models.User{
		UserID:   claims["user_id"].(string),
		FullName: claims["fullname"].(string),
		Email:    claims["email"].(string),
		Status:   claims["status"].(string),
		Picture:  claims["picture"].(string),
	}

	newAccess, err := middleware.GenerateAccessToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot generate new access token"})
	}

	return c.JSON(fiber.Map{
		"access_token": newAccess,
	})
}
