package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB // global connection

func Connect() error {
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, pass, host, port, name)

	var err error
	DB, err = sql.Open("mysql", dsn) // ใช้ global ตัวนี้
	if err != nil {
		return fmt.Errorf("cannot open DB: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("cannot ping DB: %w", err)
	}

	return nil
}
