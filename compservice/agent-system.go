package compservice

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"skko-gateway-auth/db"

	"github.com/gofiber/fiber/v2"
)

// ข้อมูลที่จะส่งไป API
type DiskSummary struct {
	Count   int     `json:"count"`
	TotalGB float64 `json:"total_gb"`
}
type SystemInfo struct {
	Idclient   int         `json:"device_id"`
	Hostname   string      `json:"hostname"`
	OS         string      `json:"os"`
	RAMTotalGB float64     `json:"ram_total_gb"`
	Disk       DiskSummary `json:"disk"`
	Timestamp  time.Time   `json:"timestamp"`
}

func AgentSystem(c *fiber.Ctx) error {
	var body SystemInfo
	if err := c.BodyParser(&body); err != nil {
		log.Printf("❌ Parse error: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request: " + err.Error(),
		})
	}

	log.Printf("📦 Received system body: %+v\n", body)

	return c.JSON(fiber.Map{
		"system": body,
	})
}

type Heartbeat struct {
	Idclient  string                 `json:"device_id"`
	Timestamp *time.Time             `json:"timestamp"`
	Meta      map[string]interface{} `json:"meta"`
}

func AgentHeartbeat(c *fiber.Ctx) error {
	var hb Heartbeat
	if err := c.BodyParser(&hb); err != nil {
		log.Printf("❌ Heartbeat parse error: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid heartbeat: " + err.Error(),
		})
	}
	metaJSON, err := json.Marshal(hb.Meta)
	if err != nil {
		log.Printf("❌ Marshal meta error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot encode meta to JSON",
		})
	}
	query := `update officedd_hardware.device set heartbeath = NOW(),meta=?,local_datetime =? WHERE device_id=?`
	res, err := db.DB.Exec(query, string(metaJSON), hb.Timestamp, hb.Idclient)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		// ไม่มีแถวที่ถูกอัปเดต
		return fmt.Errorf("device_id %v not found", hb.Idclient)
	}

	log.Printf("💓 Received heartbeat: %+v\n", hb)
	return nil

	// // update last seen client timestamp in memory or DB

	// return c.JSON(fiber.Map{
	// 	"status": hb,
	// })
}
