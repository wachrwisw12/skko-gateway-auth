package middleware

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
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
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing_token"})
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")

		// 1️⃣ Parse JWT โดยไม่ validate claims (ยังไม่เช็ก exp)
		token, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return utils.PublicKey, nil
		}, jwt.WithoutClaimsValidation())
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token_invalid"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token_invalid_claims"})
		}

		sessionID, ok := claims["sessionId"].(string)
		if !ok || sessionID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing_session_id"})
		}

		// 2️⃣ ตรวจสอบ session status ใน DB
		var status string
		var userID int
		err = db.DB.QueryRow(`SELECT user_id, status FROM user_sessions_app WHERE session_id=?`, sessionID).Scan(&userID, &status)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "session_not_found"})
		}

		if status != "active" {
			// revoked / inactive → block
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token_revoked"})
		}

		// 3️⃣ ค่อยตรวจสอบ expired
		if expVal, ok := claims["exp"].(float64); ok {
			expTime := time.Unix(int64(expVal), 0)
			if time.Now().After(expTime) {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token_expired"})
			}
		} else {
			// ถ้าไม่มี exp claim
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token_missing_exp"})
		}

		// ✅ ผ่านทุกอย่าง → เซ็ตข้อมูลไว้ใน context
		c.Locals("session_id", sessionID)
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

	// บันทึก session ใหม่
	_, err = tx.Exec(`
        INSERT INTO user_sessions_app(
            session_id, user_id,device,access_token,refresh_token,status,access_expires_at,refresh_expires_at
        ) VALUES (?, ?, ?, ?, ?, 'active', ?, ?)`, sessionId, userId, device, accessToken, refreshToken,
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
        SET access_token=?, refresh_token=?, access_expires_at=?, refresh_expires_at=? ,access_token_update_at=?
        WHERE session_id=?`,
		newAccess, newRefresh, time.Now().Add(30*time.Minute), time.Now().Add(100*24*time.Hour), time.Now().Add(100*24*time.Hour), sessionID,
	)
	return
}

// คืน true ถ้า token หมดอายุ, false ถ้ายังใช้ได้
func IsTokenExpired(token string) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return true // token ผิดรูปแบบ → ถือว่าหมดอายุ
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return true // decode ไม่ได้ → หมดอายุ
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return true // unmarshal ไม่ได้ → หมดอายุ
	}

	expVal, ok := payload["exp"]
	if !ok {
		return true // ไม่มี exp → หมดอายุ
	}

	expFloat, ok := expVal.(float64)
	if !ok {
		return true // exp ไม่ใช่ตัวเลข → หมดอายุ
	}

	return time.Now().After(time.Unix(int64(expFloat), 0))
}
