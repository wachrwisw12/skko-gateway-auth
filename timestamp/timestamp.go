package timestamp

import (
	"fmt"
	"strconv"
	"strings"
	"time"

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
			"error": err,
		})
	}
	// เรียกฟังก์ชันเช็ค leave
	rowsExist, err := HasLeave(UserID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// checkin_date_time, checkout_date_time, err := CheckTimeInOut(UserID)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 		"error": err.Error(),
	// 	})
	// }
	parts := strings.Split(locateOffice, ",") // แยกด้วย comma
	if len(parts) != 2 {
		fmt.Println("พิกัดไม่ถูกต้องs")
	}

	lat, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	lng, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)

	if err1 != nil || err2 != nil {
		fmt.Println("แปลงเป็น float ไม่สำเร็จ")
	}

	fmt.Println("Latitude:", lat)
	fmt.Println("Longitude:", lng)
	now := time.Now().In(loc)
	serverDateTime := now.Format("2006-01-02 15:04:05")
	print(rowsExist)
	return c.JSON(fiber.Map{
		"allowedLat":     lat,
		"allowedLng":     lng,
		"allowedRadius":  50,
		"serverDateTime": serverDateTime,
		"hasLeave":       rowsExist,
		// "checkin_date_time":  checkin_date_time,
		// "checkout_date_time": checkout_date_time,
	})
}

type CheckinRequest struct {
	Type string  `json:"type"` // "in" หรือ "out"
	Lat  float64 `json:"lat"`
	Lng  float64 `json:"lng"`
}

func TimestampCheckIn(c *fiber.Ctx) error {
	// อ่าน timezone Asia/Bangkok
	// loc, _ := time.LoadLocation("Asia/Bangkok")
	// now := time.Now().In(loc)

	// // parse body
	var body CheckinRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request: " + err.Error(),
		})
	}

	// // สมมติว่ามี userId มาจาก middleware (JWT)
	// userId := c.Locals("userId")
	// if userId == nil {
	// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
	// 		"error": "Unauthorized",
	// 	})
	// }

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
	return c.JSON(fiber.Map{
		"message": "Check-in success",
		// "data":    checkinData,
		"body":               body.Type,
		"checkin_date_time":  "",
		"checkout_date_time": "",
	})
}
