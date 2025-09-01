package middleware

import (
	"database/sql"
	"log"
	"strings"
	"time"

	"skko-gateway-auth/db"
	"skko-gateway-auth/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return utils.PublicKey, nil
		})
		if err != nil {
			if strings.Contains(err.Error(), "token is expired") {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token_expired"})
			}
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}

		// ตรวจสอบ session_id กับ DB
		var status string
		var userID int

		err = db.DB.QueryRow(`SELECT user_id,status FROM user_sessions_app WHERE session_id=?`, claims["sessionId"]).Scan(&userID, &status)
		if err != nil || status != "active" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "session_inactive"})
		}

		c.Locals("session_id", claims["sessionId"])
		c.Locals("user_id", userID)
		return c.Next()
	}
}

func GenerateTokensAndSaveSession(userId int, device string) (accessToken, sessionId string, err error) {
	utils.LoadKeys()
	sessionId = uuid.NewString()
	log.Println("device is", device)
	// สร้าง token
	accessClaims := jwt.MapClaims{
		"sessionId": sessionId,
		"exp":       time.Now().Add(30 * time.Minute).Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	accessToken, err = at.SignedString(utils.PrivateKey)
	if err != nil {
		return
	}

	refreshClaims := jwt.MapClaims{
		"sid": sessionId,
		"exp": time.Now().Add(100 * 24 * time.Hour).Unix(),
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	refreshToken, err := rt.SignedString(utils.PrivateKey)
	if err != nil {
		return
	}

	// Transaction
	tx, err := db.DB.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// นับ session
	var count int
	err = tx.QueryRow(`
        SELECT COUNT(*) FROM user_sessions_app 
        WHERE user_id = ? AND status = 'active'
    `, userId).Scan(&count)
	if err != nil {
		return
	}

	// // ลบ session เก่าที่สุดถ้าเกิน 3
	// if count >= 3 {
	// 	var oldestSession string
	// 	err = tx.QueryRow(`
	//         SELECT session_id FROM user_sessions_app
	//         WHERE user_id = ? AND status = 'active'
	//         ORDER BY access_expires_at ASC LIMIT 1
	//     `, userId).Scan(&oldestSession)
	// 	if err != nil {
	// 		return
	// 	}
	// 	_, err = tx.Exec(`UPDATE  user_sessions_app SET status='revoked' WHERE session_id = ?`, oldestSession)
	// 	if err != nil {
	// 		return
	// 	}
	// }

	// บันทึก session ใหม่
	_, err = tx.Exec(`
        INSERT INTO user_sessions_app(
            session_id, user_id,device,access_token,refresh_token,status,access_expires_at,refresh_expires_at
        ) VALUES (?, ?, ?, ?, ?, 'active', ?, ?)
    `, sessionId, userId, device, accessToken, refreshToken,
		time.Now().Add(30*time.Minute), time.Now().Add(100*24*time.Hour))
	return
}

// RefreshSessionToken รับ session_id แล้ว return access + refresh token ใหม่
func RefreshSessionToken(sessionID string) (newAccess, newRefresh string, err error) {
	var refreshToken string
	var status string
	var refreshExp time.Time

	// 1️⃣ Query session จาก DB
	err = db.DB.QueryRow(`
        SELECT refresh_token, status, refresh_expires_at 
        FROM user_sessions_app 
        WHERE session_id = ?
    `, sessionID).Scan(&refreshToken, &status, &refreshExp)
	if err != nil {
		if err == sql.ErrNoRows {
			err = fiber.ErrUnauthorized
		}
		return
	}

	if status != "active" {
		err = fiber.ErrUnauthorized
		return
	}

	// 2️⃣ Verify refresh token
	token, parseErr := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return utils.PublicKey, nil
	})
	if parseErr != nil || !token.Valid || time.Now().After(refreshExp) {
		err = fiber.ErrUnauthorized
		return
	}

	// 3️⃣ สร้าง access token ใหม่
	accessClaims := jwt.MapClaims{
		"sessionId": sessionID,
		"exp":       time.Now().Add(30 * time.Minute).Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims)
	newAccess, err = at.SignedString(utils.PrivateKey)
	if err != nil {
		return
	}

	// 4️⃣ สร้าง refresh token ใหม่ (rotate)
	refreshClaims := jwt.MapClaims{
		"sid": sessionID,
		"exp": time.Now().Add(100 * 24 * time.Hour).Unix(),
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	newRefresh, err = rt.SignedString(utils.PrivateKey)
	if err != nil {
		return
	}

	// 5️⃣ Update DB
	_, err = db.DB.Exec(`
        UPDATE user_sessions_app 
        SET access_token=?, refresh_token=?, access_expires_at=?, refresh_expires_at=? 
        WHERE session_id=?`,
		newAccess, newRefresh, time.Now().Add(30*time.Minute), time.Now().Add(100*24*time.Hour), sessionID,
	)
	return
}
