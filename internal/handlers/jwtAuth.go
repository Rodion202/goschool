package handlers

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

// Claims описывает содержимое нашего токена.
// Секрет передаётся явно, а не читается из окружения на инициализации пакета.
type Claims struct {
	Role string `json:"role"`
	jwt.StandardClaims
}

func GenerateJWT(secret []byte, userID int64, role string) (string, error) {
	now := time.Now()
	claims := Claims{
		Role: role,
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(72 * time.Hour).Unix(),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}
