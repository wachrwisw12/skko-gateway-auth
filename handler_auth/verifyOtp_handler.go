package handler_auth

import (
	"fmt"

	"skko-gateway-auth/db"
	"skko-gateway-auth/middleware"
	"skko-gateway-auth/models"

	"github.com/gofiber/fiber/v2"
)

func VerifyOtpHandler(c *fiber.Ctx) error {
	var body models.RequestOtp
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request json body",
		})
	}

	statusVerify := middleware.VerifyOTP(body.Email, body.Otp)
	fmt.Println(body.Email)

	if statusVerify {
		query := `
			SELECT user_id, CONCAT(prename, user_first_name, ' ', user_last_name) AS fullname
			FROM co_user
			WHERE email = ?
		`
		var user models.User // <-- struct ที่จะใส่ผลลัพธ์
		err := db.DB.QueryRow(query, body.Email).Scan(&user.UserID, &user.FullName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to fetch user",
			})
		}
		token, err := middleware.GenerateJWT(user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "ไม่สามารถสร้าง token ได้",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"verifyStatus": statusVerify,
			"token":        token,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"verifyStatus": statusVerify,
		})
	}
}
