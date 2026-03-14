package app

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash), err
}

func CheckPasswordHash(hash string, password1 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password1))
	if err != nil {
		return false
	}
	return true
}
