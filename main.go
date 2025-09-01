package main

import (
	"database/sql"
	"log"
	"os"

	"skko-gateway-auth/db"
	"skko-gateway-auth/routes"
	"skko-gateway-auth/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func main() {
	// โหลด .env
	err := godotenv.Load() // โหลด .env
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	secret := os.Getenv("SECRET_KEY_GATEWAY")
	if secret != "skko-GateWayAuth" {
		log.Fatal("❌ SECRET_KEY_GATEWAY ไม่ถูกต้อง — ปิดโปรแกรม")
	}

	log.Println("✅ รหัสผ่านถูกต้อง เริ่มระบบได้")
	// เชื่อมต่อ DB
	if err := db.Connect(); err != nil {
		log.Fatal("❌ ไม่สามารถเชื่อมต่อ database ได้:", err)
	}
	log.Println("db connected")
	utils.LoadKeys()
	app := fiber.New()
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
