package models

import "time"

type CheckinRequest struct {
	// Type string  `json:"type"` // "in" หรือ "out"
	Lat                  float64   `json:"lat"`
	Lng                  float64   `json:"lng"`
	Workingtableid       int       `json:"working_table_id"`
	CheckinDatetime      time.Time `json:"checkin_datetime"`
	Period               int       `json:"period`
	WorkingStartDateTime time.Time `json:working_start_datetime`
}
