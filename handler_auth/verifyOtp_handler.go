package handler_auth

import (
	"database/sql"

	"skko-gateway-auth/db"
	"skko-gateway-auth/middleware"
	"skko-gateway-auth/models"
	"skko-gateway-auth/services"

	"github.com/gofiber/fiber/v2"
)

func VerifyOtpHandler(c *fiber.Ctx) error {
	var body models.VerifyOTPRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request json body",
		})
	}
	var device string
	if body.Device == "" {
		device = c.Get("User-Agent")
	} else {
		device = body.Device
	}
	// ดึง device จาก header

	// ตรวจสอบ OTP
	query := `
		SELECT uuid, otp_code, insert_user_id
		FROM otp
		WHERE uuid = ? AND otp_code = ? AND expired_datetime > NOW();
	`

	var user models.VerifyOTPRequest
	err := db.DB.QueryRow(query, body.Uuid, body.OtpCode).Scan(&user.Uuid, &user.OtpCode, &user.Userid)
	if err == sql.ErrNoRows {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "รหัส OTP ไม่ถูกต้องหรือหมดเวลา",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	userInfo, err := services.GetUserInfo(user.Userid)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}
	// สร้าง token และ session
	accessToken, sessionId, err := middleware.GenerateTokensAndSaveSession(user.Userid, device)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
	}

	return c.JSON(fiber.Map{
		"message":      "OTP Verified",
		"session_id":   sessionId,
		"access_token": accessToken,
		"user_info":    userInfo,
	})
}

func CheckQrcodeHandler(c *fiber.Ctx) error {
	var body models.VerifyOTPRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request json body",
		})
	}

	// Query: คืน row พร้อม status valid/expired
	row := db.DB.QueryRow(`
        SELECT uuid,
               CASE
                   WHEN expired_datetime > NOW() THEN 'valid'
                   ELSE 'expired'
               END AS status
        FROM otp
        WHERE uuid = ?;
    `, body.Uuid)

	var id string
	var status string
	err := row.Scan(&id, &status)
	if err != nil {
		if err == sql.ErrNoRows {
			// OTP ไม่เจอ
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "QRCODE ไม่ถูกต้อง",
			})
		}
		// error จริงจาก DB
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// เช็ก status ที่ query คืนมา
	if status == "expired" {
		return c.Status(fiber.StatusGone).JSON(fiber.Map{
			"error": "OTP expired",
		})
	}

	return c.JSON(fiber.Map{
		"message": "QRCODE valid",
	})
}
