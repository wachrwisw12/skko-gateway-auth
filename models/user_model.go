package models

type User struct {
	FullName       string `json:"fullname"`
	Picture        string `json:"picture"`
	SexID          int    `json:"sex_id"`
	MainOfficeID   int    `json:"main_office_id"`
	MainOfficeName string `json:"main_office_name"`
	SubOfficeID    int    `json:" sub_office_id"`
	SubOfficeName  string `json:"sub_office_name"`
}
