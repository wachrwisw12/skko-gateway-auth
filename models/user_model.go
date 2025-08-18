package models

type User struct {
	UserID   string `json:"user_id"`
	FullName string `json:"fullname"`
	Status   string `json:"status"`
	Email    string `json:"email"`
	Picture  string `json:"picture"`
}
