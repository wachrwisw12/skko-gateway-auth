package models

type LeaveInfo struct {
	UserID    int    `json:"user_id"`
	UserName  string `json:"user_name"`
	LeaveType string `json:"leave_type"`
	LeaveDate string `json:"leave_date"`
}
