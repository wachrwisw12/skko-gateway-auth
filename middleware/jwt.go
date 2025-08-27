package middleware

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"skko-gateway-auth/models"

	"github.com/golang-jwt/jwt/v5"

	"github.com/gofiber/fiber/v2"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
		}

		token := strings.TrimPrefix(auth, "Bearer ")

		claims, err := ParseJWT(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}

		c.Locals("user_id", claims["user_id"])
		return c.Next()
	}
}

func GenerateAccessToken(user models.User) (string, error) {
	return generateToken(user, 30*time.Minute) // อายุ 5 นาที
}

func GenerateRefreshToken(user models.User) (string, error) {
	return generateToken(user, 7*24*time.Hour) // อายุ 7 วัน
}

func generateToken(user models.User, duration time.Duration) (string, error) {
	// แปลง struct → map[string]interface{} ด้วย JSON
	userMap := make(map[string]interface{})

	b, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("marshal user failed: %w", err)
	}

	if err := json.Unmarshal(b, &userMap); err != nil {
		return "", fmt.Errorf("unmarshal user failed: %w", err)
	}

	// Claims
	claims := jwt.MapClaims{
		// "user_id": user.UserID,
		// "user":    userMap,
		"exp": time.Now().Add(duration).Unix(),
	}
	// สร้าง token ด้วย HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// อ่าน JWT_SECRET จาก environment
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", fmt.Errorf("missing JWT_SECRET environment variable")
	}

	// สร้าง JWT string
	return token.SignedString([]byte(secret))
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
