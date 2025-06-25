package token

import (
	"time"

	"github.com/Hiendang123/golang-server.git/internal/common"
	"github.com/golang-jwt/jwt"
)

func GenerateToken(userID uint, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(common.SecretKey))
}
