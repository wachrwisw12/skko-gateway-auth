package models

import "time"

type OtpEntry struct {
	Code      string
	ExpiresAt time.Time
}

type RequestOtp struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}
