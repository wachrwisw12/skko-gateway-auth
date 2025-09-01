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
	})
}
