package handler_auth

import (
	"database/sql"
	"fmt"

	"skko-gateway-auth/db"
	"skko-gateway-auth/models"

	"github.com/gofiber/fiber/v2"
)

func VerifyOtpHandler(c *fiber.Ctx) error {
	var body models.VerifyOTPRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request json body",
		})
	}

	// statusVerify := middleware.VerifyOTP(body.Uuid)
	fmt.Println(body.Uuid)

	query := `
	 		SELECT uuid ,otp_code
			FROM otp
			WHERE uuid = ? AND otp_code=? AND expired_datetime > NOW(); 
 	`
	var user models.VerifyOTPRequest // <-- struct ที่จะใส่ผลลัพธ์
	err := db.DB.QueryRow(query, body.Uuid, body.OtpCode).Scan(&user.Uuid, &user.OtpCode)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "รหัส OTP ไม่ถูกต้อง",
		})
	} else if err != nil {
		return err
	}

	return nil
}

func CheckQrcodeHandler(c *fiber.Ctx) error {
	var body models.VerifyOTPRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request json body",
		})
	}

	// statusVerify := middleware.VerifyOTP(body.Uuid)
	fmt.Println(body.Uuid)

	query := `
	 		SELECT uuid 
			FROM otp
			WHERE uuid = ?  AND expired_datetime > NOW(); 
 	`
	var user models.VerifyOTPRequest // <-- struct ที่จะใส่ผลลัพธ์
	err := db.DB.QueryRow(query, body.Uuid).Scan(&user.Uuid)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "QRCODE หมดเวลา",
		})
	} else if err != nil {
		return err
	}

	return nil
}
