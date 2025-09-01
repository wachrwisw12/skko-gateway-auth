package services

import (
	"database/sql"

	"skko-gateway-auth/db"
	"skko-gateway-auth/models"
)

func GetUserInfo(user_id int) (*models.User, error) {
	query := `
        SELECT  CONCAT(u.prename,u.user_first_name,'  ',u.user_last_name) as fullname,u.sex_id,u.picture FROM co_user u
        WHERE u.user_id=? 
    `
	row := db.DB.QueryRow(query, user_id)

	var user models.User

	err := row.Scan(&user.FullName, &user.SexID, &user.Picture)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &user, nil
}
