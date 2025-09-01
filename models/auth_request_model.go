package models

import "time"

type OtpEntry struct {
	Code      string
	ExpiresAt time.Time
}

type VerifyOTPRequest struct {
	Uuid    string `json:"uuid"`
	OtpCode string `json:"otp_code"`
	Userid  int    `json:"user_id"`
	Device  string `json:"device"`
}
