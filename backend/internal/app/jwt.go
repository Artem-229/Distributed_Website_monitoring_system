package app

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateJWTToken(idd uuid.UUID, secretword string) (string, error) {
	token := jwt.MapClaims{
		"userid": idd,
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}

	newtoken := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
	return newtoken.SignedString([]byte(secretword))
}
