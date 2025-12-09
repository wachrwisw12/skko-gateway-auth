package services

import (
	"database/sql"

	"skko-gateway-auth/db"
	"skko-gateway-auth/models"
)

func GetUserInfo(user_id int) (*models.User, error) {
	query := `
SELECT 
    CONCAT(u.prename, u.user_first_name, ' ', u.user_last_name) AS fullname,
    u.sex_id,
    u.picture,
    ofe.main_office_id,
    main_office.office_name_2 AS main_office_name,
    SUBSTRING_INDEX(SUBSTRING_INDEX(ofe.office_relation_code, '-', 2), '-', -1) AS sub_office_id,
    sub_office.office_name AS sub_office_name
FROM co_user u
LEFT JOIN co_office ofe 
    ON ofe.office_id = u.office_id
LEFT JOIN co_office main_office 
    ON main_office.office_id = ofe.main_office_id
LEFT JOIN co_office sub_office
    ON sub_office.office_id = CAST(SUBSTRING_INDEX(SUBSTRING_INDEX(ofe.office_relation_code, '-', 2), '-', -1) AS UNSIGNED)
WHERE u.user_id = ?
`

	row := db.DB.QueryRow(query, user_id)

	var user models.User

	err := row.Scan(&user.FullName, &user.SexID, &user.Picture, &user.MainOfficeID, &user.MainOfficeName, &user.SubOfficeID, &user.SubOfficeName)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
