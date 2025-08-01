package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"github.com/wachrwisw12/corework-gateway-auth/db"
	"github.com/wachrwisw12/corework-gateway-auth/routes"
)

var DB *sql.DB

func main() {
	// โหลด .env
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ ไม่พบไฟล์ .env (ไม่เป็นไรถ้า set env ไว้ในระบบแล้ว)")
	}
	secret := os.Getenv("SECRET_KEY_GATEWAY")
	if secret != "airhosGateWayAuth" {
		log.Fatal("❌ SECRET_KEY_GATEWAY ไม่ถูกต้อง — ปิดโปรแกรม")
	}

	log.Println("✅ รหัสผ่านถูกต้อง เริ่มระบบได้")
	// เชื่อมต่อ DB
	if err := db.Connect(); err != nil {
		log.Fatal("❌ ไม่สามารถเชื่อมต่อ database ได้:", err)
	}

	app := fiber.New()
	routes.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
