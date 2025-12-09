package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB // global connection

func Connect() error {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	// DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Asia%%2FBangkok&timeout=5s&readTimeout=30s&writeTimeout=30s",
		user, pass, host, port, name)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("cannot open DB: %w", err)
	}

	// ปรับ connection pool ให้จัดการ connection ได้ดีขึ้น
	DB.SetMaxOpenConns(20)                 // จำนวน connection สูงสุด
	DB.SetMaxIdleConns(10)                 // idle connection ที่เก็บไว้
	DB.SetConnMaxLifetime(5 * time.Minute) // อายุ connection ไม่เกินนี้ → driver จะสร้างใหม่เอง
	DB.SetConnMaxIdleTime(2 * time.Minute) // idle connection เกินนี้จะถูกปิด

	// ทดสอบเชื่อมต่อ
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("cannot ping DB: %w", err)
	}

	return nil
}
