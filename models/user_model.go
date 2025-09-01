package models

type User struct {
	FullName string `json:"fullname"`
	Picture  string `json:"picture"`
	SexID    int    `json:"sex_id"`
}
