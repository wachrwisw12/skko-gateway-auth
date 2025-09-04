package timestamp

import (
	"database/sql"

	"skko-gateway-auth/db"
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
	query := `SELECT 
        COALESCE(
            NULLIF(offi.gps, ''),  
            fallback.gps
        ) AS gps
    FROM co_user usr
    LEFT JOIN co_office offi
        ON offi.office_id = usr.office_id
    LEFT JOIN co_office fallback
        ON fallback.office_id = CAST(SUBSTRING_INDEX(offi.office_relation_code, '-', 1) AS UNSIGNED)
    WHERE usr.user_id = ?;`

	var gps string
	err := db.DB.QueryRow(query, UserID).Scan(&gps)
	if err != nil {
		return "", err
	}

	return gps, nil
}

// func CheckTimeInOut(UserID int) (string, string, error) {
// 	query := `SELECT checkin_date_time,checkout_date_time FROM timestamp WHERE user_id=? AND date = DATE(NOW())`
// 	var checkin_date_time string
// 	var checkout_date_time string
// 	err := db.DB.QueryRow(query, UserID).Scan(&checkin_date_time, &checkout_date_time)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	return checkin_date_time, checkout_date_time, nil
// }
