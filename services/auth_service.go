package services

import (
	"database/sql"
	"errors"

	"skko-gateway-auth/db"
	"skko-gateway-auth/models"

	"golang.org/x/crypto/bcrypt"
)

func LoginByUser(username string, password string) (*models.User, error) {
	query := `
        SELECT u.id, CONCAT(u.name) as fullname, u.status, u.password
        FROM users u
       
        WHERE u.username=? 
    `
	row := db.DB.QueryRow(query, username)

	var user models.User
	var hashedPassword string

	err := row.Scan(&user.UserID, &user.FullName, &user.Status, &hashedPassword)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	// เปรียบเทียบ password plain กับ bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, errors.New("รหัสผ่านไม่ถูกต้องd")
	}

	return &user, nil
}

func LoginByEmail(email string, password string) (*models.User, error) {
	query := `
        SELECT u.user_id, CONCAT(u.prename,u.user_first_name,' ',u.user_last_name) as fullname,u.picture,u.password
        FROM co_user u
    
        WHERE u.email=? 
    `
	row := db.DB.QueryRow(query, email)

	var user models.User
	var hashedPassword string

	err := row.Scan(&user.UserID, &user.FullName, &user.Picture, &hashedPassword)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	// เปรียบเทียบ password plain กับ bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, err
	}
	return &user, nil
}
