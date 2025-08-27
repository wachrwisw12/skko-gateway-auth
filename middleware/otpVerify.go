package middleware

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	"skko-gateway-auth/models"

	"gopkg.in/gomail.v2"
)

var (
	otpStore = make(map[string]models.OtpEntry)
	otpMutex sync.Mutex
)

func GenerateOTP(email string) (string, error) {
	const digits = "0123456789"
	otp := ""
	for i := 0; i < 6; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		otp += string(digits[num.Int64()])
	}

	otpMutex.Lock()
	defer otpMutex.Unlock()

	otpStore[email] = models.OtpEntry{
		Code:      otp,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	return otp, nil
}

func VerifyOTP(email, otp string) bool {
	otpMutex.Lock()
	defer otpMutex.Unlock()

	entry, exists := otpStore[email]
	if !exists {
		return false
	}

	if time.Now().After(entry.ExpiresAt) {
		delete(otpStore, email) // ลบ OTP ที่หมดอายุแล้ว
		return false
	}

	if entry.Code != otp {
		return false
	}

	// ตรวจสอบสำเร็จ — ลบ OTP ทิ้งเพื่อความปลอดภัย
	delete(otpStore, email)
	return true
}

func SendOTPEmail(to, otp string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "ict3.sakonnakhon@gmail.com") // เปลี่ยนเป็นอีเมลคุณ
	m.SetHeader("To", to)
	m.SetHeader("Subject", "รหัส OTP ของคุณ")
	m.SetBody("text/plain", fmt.Sprintf("รหัส OTP ของคุณคือ: %s (หมดอายุภายใน 5 นาที)", otp))

	d := gomail.NewDialer("smtp.gmail.com", 587, "ict3.sakonnakhon@gmail.com", os.Getenv("APP_GOOGLE_PASS"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
