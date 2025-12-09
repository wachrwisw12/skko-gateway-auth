package timestamp

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"skko-gateway-auth/db"
	"skko-gateway-auth/models"
)

func HasLeave(UserID int) (bool, error) {
	query := `SELECT user_id,leave_id FROM officedd_timestamp.leave WHERE user_id=? AND date = DATE(NOW())`
	var dummy int
	err := db.DB.QueryRow(query, UserID).Scan(&dummy)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetLocalOffice(UserID int) (string, error) {
	// ดึง office ของ user
	query := `
        SELECT offi.gps, offi.belong_to_office_id, offi.office_id, offi.main_office_id
        FROM co_office offi
        LEFT JOIN co_user co ON co.office_id = offi.office_id
        WHERE co.user_id = ?;
    `

	var gps sql.NullString
	var belongTo sql.NullString
	var officeID string
	var mainOfficeID sql.NullString

	err := db.DB.QueryRow(query, UserID).Scan(&gps, &belongTo, &officeID, &mainOfficeID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}

	// ถ้า office มี gps → return เลย
	if gps.Valid && gps.String != "" {
		return gps.String, nil
	}

	// ถ้าไม่มี gps → ไล่ parent belong_to_office_id
	currentID := belongTo.String
	visited := map[string]bool{} // กัน loop วนซ้ำ

	for currentID != "" {
		if visited[currentID] {
			break
		}
		visited[currentID] = true

		var parentGPS sql.NullString
		var parentBelong sql.NullString

		err := db.DB.QueryRow(
			`SELECT gps, belong_to_office_id FROM co_office WHERE office_id = ?`,
			currentID,
		).Scan(&parentGPS, &parentBelong)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				break
			}
			return "", err
		}

		if parentGPS.Valid && parentGPS.String != "" {
			return parentGPS.String, nil
		}

		currentID = parentBelong.String
	}

	// ถ้ายังไม่เจอ → fallback ไป main_office_id
	if mainOfficeID.Valid && mainOfficeID.String != "" {
		var mainGPS sql.NullString
		err := db.DB.QueryRow(
			`SELECT gps FROM co_office WHERE office_id = ?`,
			mainOfficeID.String,
		).Scan(&mainGPS)

		if err == nil && mainGPS.Valid && mainGPS.String != "" {
			return mainGPS.String, nil
		}
	}

	// ถ้าไม่เจอ gps เลย → return default
	return "", nil
}

type StautCheckin struct {
	WorkingTableId                int        `json:"working_table_id"`
	Period                        int        `json:"period"`
	LeaveLetterID                 *int       `json:"leave_letter_id"`
	LeaveCancelID                 *int       `json:"leave_cancel_id"`
	CheckinDateTime               *time.Time `json:"checkin_datetime"`
	CheckoutDateTime              *time.Time `json:"checkout_datetime"`
	TimestampStartableMinDatetime *time.Time `json:"timestamp_startable_min_datetime"`
	TimestampEndableMaxDatetime   *time.Time `json:"timestamp_endable_max_datetime"`
	TimestampEndableMinDatetime   *time.Time `json:"timestamp_endable_min_datetime"`
	WorkingStartDateTime          *time.Time `json:"working_start_datetime"`
}
type TimeStampHistorys struct {
	WorkingDate     *time.Time `json:"working_date"`
	CheckinDatetime *time.Time `json:"checkin_datetime"`
	CheckinResult   *string    `json:"checkin_result"`
}

func CheckinState(UserID int) ([]StautCheckin, error) {
	query := `
        SELECT wt.working_table_id,wt.period_id, wt.leave_letter_id,wt.leave_cancel_id,wt.checkin_datetime,wt.checkout_datetime,wt.timestamp_startable_min_datetime,wt.timestamp_endable_max_datetime,wt.timestamp_endable_min_datetime,wt.working_start_datetime
        FROM officedd_timestamp.working_table wt
        WHERE wt.user_id = ? AND wt.working_date = CURRENT_DATE()`
	rows, err := db.DB.Query(query, UserID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var working []StautCheckin
	for rows.Next() {
		var p StautCheckin
		if err := rows.Scan(&p.WorkingTableId, &p.Period, &p.LeaveLetterID, &p.LeaveCancelID, &p.CheckinDateTime, &p.CheckoutDateTime, &p.TimestampStartableMinDatetime, &p.TimestampEndableMaxDatetime, &p.TimestampEndableMinDatetime, &p.WorkingStartDateTime); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}
		working = append(working, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	// ถ้าไม่มีข้อมูล ให้คืนค่า default object

	return working, nil
}

func Checkin(UserID int, body models.CheckinRequest) error {
	// เรียกฟังก์ชันเช็คสาย
	status, err := CheckLate(body.CheckinDatetime, body.WorkingStartDateTime)
	if err != nil {
		return err
	}
	queryCheckinin := `update officedd_timestamp.working_table set checkin_datetime = NOW(), checkin_lat=? ,checkin_lng=?,checkin_result=? WHERE working_table_id= ?`
	res, err := db.DB.Exec(queryCheckinin, body.Lat, body.Lng, status, body.Workingtableid)
	if err != nil {
		fmt.Print(res)
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		fmt.Print(res)
		return err
	}
	if row == 0 {
		fmt.Print(res)
		return fmt.Errorf("ลงเวลาเข้าไม่สำเร็จ")
	}
	fmt.Print(res)
	return nil
}

func Checkout(UserID int, body models.CheckinRequest) error {
	queryCheckinout := `update officedd_timestamp.working_table set checkout_datetime = NOW(), checkout_lat=? ,checkout_lng=? WHERE working_table_id= ?`
	res, err := db.DB.Exec(queryCheckinout, body.Lat, body.Lng, body.Workingtableid)
	if err != nil {
		return err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return fmt.Errorf("ลงเวลาออกไม่สำเร็จ")
	}

	return nil
}

func GetTimestampHistory(UserID int) ([]TimeStampHistorys, error) {
	query := `SELECT working_date,checkin_datetime,checkin_result FROM officedd_timestamp.working_table ots WHERE ots.user_id=? AND working_date BETWEEN DATE_FORMAT(CURDATE(), '%Y-%m-01') AND CURRENT_DATE()`
	rows, err := db.DB.Query(query, UserID)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()
	var historyData []TimeStampHistorys
	for rows.Next() {
		var p TimeStampHistorys
		if err := rows.Scan(&p.WorkingDate, &p.CheckinDatetime, &p.CheckinResult); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		historyData = append(historyData, p)
	}
	return historyData, nil
}

func CheckLate(CheckinDatetime time.Time, WorkingStartDatetime time.Time) (string, error) {
	// ดึง timezone
	location := CheckinDatetime.Location()

	// ดึงวัน เดือน ปี ของวันที่เช็กอิน
	y, m, d := CheckinDatetime.Date()

	// กำหนดเวลาเริ่มงานของวันนั้น (ใช้เวลาจาก WorkingStartDatetime เช่น 08:30)
	startHour := WorkingStartDatetime.Hour()
	startMinute := WorkingStartDatetime.Minute()

	// สร้างเวลาเริ่มงานของวันเดียวกัน
	workStart := time.Date(y, m, d, startHour, startMinute, 0, 0, location)

	// คำนวณส่วนต่างระหว่างเวลาที่เช็กอินกับเวลาเริ่มงาน
	diff := CheckinDatetime.Sub(workStart)

	// ตรวจสอบว่าสายเกิน 30 นาทีหรือไม่
	if diff > 30*time.Minute {
		return "สาย", nil
	}

	// ถ้าเช็กอินก่อนเวลาหรือภายใน 30 นาที
	return "ทันเวลา", nil
}
