package timestamp

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"skko-gateway-auth/models"

	"github.com/gofiber/fiber/v2"
)

func TimestampHome(c *fiber.Ctx) error {
	UserID := c.Locals("user_id").(int)
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil
	}
	locateOffice, err := GetLocalOffice(UserID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	parts := strings.Split(locateOffice, ",") // แยกด้วย comma
	if len(parts) != 2 {
		fmt.Println("พิกัดไม่ถูกต้อง")
	}

	lat, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	lng, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)

	if err1 != nil || err2 != nil {
		fmt.Println("แปลงเป็น float ไม่สำเร็จ")
	}
	checkState, err := CheckinState(UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	now := time.Now().In(loc)
	serverDateTime := now.Format(time.RFC3339) // => จะได้ +07:00

	return c.JSON(fiber.Map{
		"allowedLat":     lat,
		"allowedLng":     lng,
		"allowedRadius":  50,
		"serverDateTime": serverDateTime, // เวลาปัจจุบัน +07:00
		"checkinState":   checkState,
	})
}

func TimestampCheckIn(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(int)

	// // parse body
	var body models.CheckinRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request: " + err.Error(),
		})
	}
	if body.Period == 1 {
		err := Checkin(userId, body)
		if err != nil {
			return err
		}
	} else if body.Period == 2 {
		err := Checkout(userId, body)
		if err != nil {
			return err
		}
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "ข้อมูลไม่ถูกต้อง",
		})
	}

	// print(userId)
	// // TODO: insert ลง DB จริง เช่น Postgres/MySQL
	// // ตัวอย่าง record
	// checkinData := fiber.Map{
	// 	"userId":    userId,
	// 	"type":      body.Type, // "in" หรือ "out"
	// 	"lat":       body.Lat,
	// 	"lng":       body.Lng,
	// 	"timestamp": now.Format("2006-01-02 15:04:05"),
	// }

	// ส่ง response กลับ
	// return c.JSON(fiber.Map{
	// 	"message":          "Check-in success",
	// 	"lat":              body.Lat,
	// 	"lng":              body.Lng,
	// 	"working_table_id": body.Workingtableid,
	// })
	return nil
}

func TimeStampHistory(c *fiber.Ctx) error {
	UserID := c.Locals("user_id").(int)
	historyData, err := GetTimestampHistory(UserID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"message": "test",
		"data":    historyData,
	})
}

// func TimeTest(c *fiber.Ctx) error {
// 	var body models.CheckinRequest
// 	if err := c.BodyParser(&body); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"error": "Invalid request: " + err.Error(),
// 		})
// 	}

// 	// เรียกฟังก์ชันเช็คสาย
// 	status, err := CheckLate(body.CheckinDatetime)
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.JSON(fiber.Map{
// 		"result": "ok",
// 		"status": status, // สาย / ทันเวลา
// 	})
// }
